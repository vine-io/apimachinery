// MIT License
//
// Copyright (c) 2021 Lack
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package rest

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"path"
	"reflect"
	"strings"
	"sync"
	"time"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"github.com/oxtoacart/bpool"
	log "github.com/vine-io/vine/lib/logger"
	"github.com/vine-io/vine/util/qson"
	"go.uber.org/atomic"
	"golang.org/x/net/netutil"
)

var (
	_ Rest = (*httpRest)(nil)

	bufferPool = bpool.NewSizedBufferPool(1024, 8)
)

type httpRest struct {
	sync.RWMutex

	serve  *http.Server
	router *gin.Engine

	exit chan chan error
	wg   *sync.WaitGroup

	options Options

	handlers map[string]Handler

	// marks the serve as started
	started *atomic.Bool
}

func newHttpRest(opts ...Option) *httpRest {
	options := NewOptions(opts...)

	//gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	hr := &httpRest{
		serve:    &http.Server{},
		router:   router,
		exit:     make(chan chan error),
		wg:       &sync.WaitGroup{},
		options:  options,
		handlers: map[string]Handler{},
		started:  atomic.NewBool(false),
	}

	hr.configure()

	return hr
}

func (h *httpRest) configure(opts ...Option) {
	for _, o := range opts {
		o(&h.options)
	}

}

func (h *httpRest) Init(opts ...Option) error {
	h.configure(opts...)
	return nil
}

func (h *httpRest) Options() Options {
	h.RLock()
	defer h.RUnlock()
	return h.options
}

func (h *httpRest) Handle(handler Handler) error {

	var err error
	opts := handler.Options()
	switch opts.Kind {
	case Resource:
		err = h.registerResourceHandler(handler.Handler(), opts)
	case Http:

	case RPC:

	}

	if err != nil {
		return err
	}

	h.handlers[handler.Name()] = handler
	return nil
}

func (h *httpRest) registerResourceHandler(handler interface{}, opts HandlerOptions) error {

	rh := handler.(ResourceHandler)
	gvk := opts.Target.Gvk
	typ := opts.Target.Target

	prefix := path.Join("/api", gvk.Group, gvk.Version, strings.ToLower(gvk.Kind))
	h.router.GET(prefix, listResourceHandler(rh, typ))
	h.router.POST(prefix, postResourceHandler(rh))
	h.router.GET(prefix+"/:uid", getResourceHandler(rh))
	h.router.PATCH(prefix+"/:uid", patchResourceHandler(rh))
	h.router.DELETE(prefix+"/:uid", deleteResourceHandler(rh))

	return nil
}

func listResourceHandler(rh ResourceHandler, typ reflect.Type) gin.HandlerFunc {
	return func(c *gin.Context) {

		_ = typ

		br, err := requestPayload(c.Request)
		if err != nil {
			c.JSON(500, fmt.Sprintf("payload: %v", err))
			return
		}

		req := &ListRequest{}
		if err = json.Unmarshal(br, req); err != nil {
			c.JSON(500, fmt.Sprintf("unmarshal %v", err))
			return
		}

		rsp := &ListResponse{}

		err = rh.List(c, req, rsp)
		if err != nil {
			return
		}

		c.JSON(200, gin.H{
			"result": gin.H{
				"list":  rsp.Out,
				"total": rsp.Total,
			},
		})
	}
}

func postResourceHandler(rh ResourceHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func getResourceHandler(rh ResourceHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func patchResourceHandler(rh ResourceHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func deleteResourceHandler(rh ResourceHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

// requestPayload takes a *http.Request.
// If the request is a GET the query string parameters are extracted and marshaled to JSON and the raw bytes are returned.
// If the request method is a POST the request body is read and returned
func requestPayload(r *http.Request) ([]byte, error) {
	var err error

	ct := r.Header.Get("Content-Type")

	// allocate maximum
	matches := make(map[string]interface{})
	bodydst := ""

	// map of all fields
	req := make(map[string]interface{})

	// get fields from url values
	if len(r.URL.RawQuery) > 0 {
		umd := make(map[string]interface{})
		err = qson.Unmarshal(&umd, r.URL.RawQuery)
		if err != nil {
			return nil, err
		}
		for k, v := range umd {
			matches[k] = v
		}
	}

	for k, v := range matches {
		ps := strings.Split(k, ".")
		if len(ps) == 1 {
			req[k] = v
			continue
		}
		em := make(map[string]interface{})
		em[ps[len(ps)-1]] = v
		for i := len(ps) - 2; i > 0; i-- {
			nm := make(map[string]interface{})
			nm[ps[i]] = em
			em = nm
		}
		if vm, ok := req[ps[0]]; ok {
			// nested map
			nm := vm.(map[string]interface{})
			for vk, vv := range em {
				nm[vk] = vv
			}
			req[ps[0]] = nm
		} else {
			req[ps[0]] = em
		}
	}
	pathbuf := []byte("{}")
	if len(req) > 0 {
		pathbuf, err = json.Marshal(req)
		if err != nil {
			return nil, err
		}
	}

	urlbuf := []byte("{}")
	out, err := jsonpatch.MergeMergePatches(urlbuf, pathbuf)
	if err != nil {
		return nil, err
	}

	switch r.Method {
	case "GET":
		// empty response
		if strings.Contains(ct, "application/json") && string(out) == "{}" {
			return out, nil
		} else if string(out) == "{}" && !strings.Contains(ct, "application/json") {
			return []byte{}, nil
		}
		return out, nil
	case "PATCH", "POST", "PUT", "DELETE":
		bodybuf := []byte("{}")
		buf := bufferPool.Get()
		defer bufferPool.Put(buf)
		if _, err := buf.ReadFrom(r.Body); err != nil {
			return nil, err
		}
		if b := buf.Bytes(); len(b) > 0 {
			bodybuf = b
		}

		if out, err = jsonpatch.MergeMergePatches(out, bodybuf); err == nil {
			return out, nil
		}

		var jsonbody map[string]interface{}
		if json.Valid(bodybuf) {
			if err = json.Unmarshal(bodybuf, &jsonbody); err != nil {
				return nil, err
			}
		}
		dstmap := make(map[string]interface{})
		ps := strings.Split(bodydst, ".")
		if len(ps) == 1 {
			if jsonbody != nil {
				dstmap[ps[0]] = jsonbody
			} else {
				// old unexpected behaviour
				dstmap[ps[0]] = bodybuf
			}
		} else {
			em := make(map[string]interface{})
			if jsonbody != nil {
				em[ps[len(ps)-1]] = jsonbody
			} else {
				// old unexpected behaviour
				em[ps[len(ps)-1]] = bodybuf
			}
			for i := len(ps) - 2; i > 0; i-- {
				nm := make(map[string]interface{})
				nm[ps[i]] = em
				em = nm
			}
			dstmap[ps[0]] = em
		}

		bodyout, err := json.Marshal(dstmap)
		if err != nil {
			return nil, err
		}

		if out, err = jsonpatch.MergeMergePatches(out, bodyout); err == nil {
			return out, nil
		}

		//fallback to previous unknown behaviour
		return bodybuf, nil
	}

	return []byte{}, nil
}

func (h *httpRest) NewHandler(handler interface{}, opts ...HandlerOption) Handler {
	return newHttpHandler(handler, opts...)
}

func (h *httpRest) Register() error {
	h.RLock()
	config := h.Options()
	h.RUnlock()

	_ = config

	// TODO: rest register

	return nil
}

func (h *httpRest) Start() error {
	if h.started.Load() {
		return nil
	}

	config := h.Options()

	var ts net.Listener
	var err error
	if tc := config.TlSConfig; tc != nil {
		ts, err = tls.Listen("tcp", config.Address, tc)
	} else {
		ts, err = net.Listen("tcp", config.Address)
	}

	if err != nil {
		return err
	}

	if config.MaxConn != 0 {
		ts = netutil.LimitListener(ts, config.MaxConn)
	}

	log.Infof("Server [http] Listening on %s", ts.Addr().String())

	h.Lock()
	h.options.Address = ts.Addr().String()
	h.Unlock()

	// announce self to the world
	if err = h.Register(); err != nil {
		log.Errorf("Rest register error: %v", err)
	}

	go func() {
		h.serve = &http.Server{
			Handler: h.router,
		}

		if err = h.serve.Serve(ts); err != nil {
			log.Errorf("Rest Server start error: %v", err)
		}
	}()

	go func() {
		t := new(time.Ticker)

		// only process if it exists
		if h.options.RegisterInterval > time.Duration(0) {
			// new ticker
			t = time.NewTicker(h.options.RegisterInterval)
		}

		// return error chan
		var ch chan error

	Loop:
		for {
			select {
			case <-t.C:
				// TODO: register
			case ch = <-h.exit:
				break Loop
			}
		}

		// deregister self
		// TODO: deregister

		// wait for waitgroup
		if h.wg != nil {
			h.wg.Wait()
		}

		// stop the http server
		exit := make(chan bool)

		ctx, cancel := context.WithTimeout(h.options.Context, time.Second)
		defer cancel()
		go func() {
			if e := h.serve.Shutdown(ctx); e != nil {
				log.Errorf("Stop Rest Server: %v", err)
			}
			close(exit)
		}()

		select {
		case <-exit:
		case <-ctx.Done():
		}

		// close transport
		ch <- nil

	}()

	// mark the server as started
	h.started.Store(true)

	return nil
}

func (h *httpRest) Stop() error {
	if !h.started.Load() {
		return nil
	}

	ch := make(chan error)
	h.exit <- ch

	var err error
	select {
	case err = <-ch:
		h.started.Store(true)
	}

	return err
}

func (h *httpRest) String() string {
	return "http"
}

func NewHttpRest(opts ...Option) *httpRest {
	return newHttpRest(opts...)
}

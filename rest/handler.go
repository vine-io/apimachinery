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
	"fmt"
	"reflect"

	"github.com/vine-io/vine/core/registry"
)

type httpHandler struct {
	name    string
	handler interface{}
	opts    HandlerOptions
}

func newHttpHandler(handler interface{}, opts ...HandlerOption) Handler {
	options := HandlerOptions{
		Kind:    Resource,
		Paths:   map[string]*registry.OpenAPIPath{},
		Models:  map[string]*registry.Model{},
		Schemes: &registry.SecuritySchemes{},
	}

	for _, o := range opts {
		o(&options)
	}

	typ := reflect.TypeOf(handler)
	hdlr := reflect.ValueOf(handler)
	name := reflect.Indirect(hdlr).Type().Name()

	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		fmt.Println(method)
	}

	return &httpHandler{
		name:    name,
		handler: handler,
		opts:    options,
	}
}

func (h *httpHandler) Name() string {
	return h.name
}

func (h *httpHandler) Handler() interface{} {
	return h.handler
}

func (h *httpHandler) Endpoints() []*registry.OpenAPIPath {
	return nil
}

func (h *httpHandler) Options() HandlerOptions {
	return h.opts
}

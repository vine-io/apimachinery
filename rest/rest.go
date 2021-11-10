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
	"time"

	"github.com/vine-io/apimachinery/runtime"
	"github.com/vine-io/vine/core/registry"
)

var (
	DefaultRest             Rest = NewHttpRest()
	DefaultAddress               = ":8080"
	DefaultRegisterInterval      = time.Second * 30
)

// Rest is a simple vine restful abstraction
type Rest interface {
	// Init initialises options
	Init(opts ...Option) error
	// Options retrieve the options
	Options() Options
	// Handle register a handler
	Handle(Handler) error
	// NewHandler create a new handler
	NewHandler(interface{}, ...HandlerOption) Handler
	// Start the rest
	Start() error
	// Stop the server
	Stop() error
	// String implementation
	String() string
}

// Handler interface represents a request handler. It's generated
// by passing any type of public concrete object with endpoints into rest.NewHandler.
// Most will pass in a struct.
//
// Example:
//
//		type Greeter struct{}
//
//		func (g *Greeter) Hello(context, request, response) error {
//			return nil
//		}
//
type Handler interface {
	Name() string
	Handler() interface{}
	Endpoints() []*registry.OpenAPIPath
	Options() HandlerOptions
}

type ResourceHandler interface {
	Get(context.Context, *GetRequest, *GetResponse) error
	List(context.Context, *ListRequest, *ListResponse) error
	Post(context.Context, *PostRequest, *PostResponse) error
	Patch(context.Context, *PatchRequest, *PatchResponse) error
	Delete(context.Context, *DeleteRequest, *DeleteResponse) error
}

type GetRequest struct {
	UID string `json:"uid"`
}

type GetResponse struct {
	Obj runtime.Object
}

type ListRequest struct {
	Page int32          `json:"page"`
	Size int32          `json:"size"`
	Ps   runtime.Object `json:"ps"`
}

type ListResponse struct {
	Out   interface{}
	Total int64
}

type PostRequest struct {
	Obj runtime.Object
}

type PostResponse struct {
	Out runtime.Object
}

type PatchRequest struct {
	UID string `json:"uid"`
	Obj runtime.Object
}

type PatchResponse struct {
	Out runtime.Object
}

type DeleteRequest struct {
	UID  string
	Soft bool
}

type DeleteResponse struct{}

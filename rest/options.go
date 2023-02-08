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
	"reflect"
	"time"

	"github.com/vine-io/apimachinery/runtime"
	"github.com/vine-io/apimachinery/schema"
	pb "github.com/vine-io/vine/lib/api/handler/openapi/proto"
)

type Options struct {
	Address string

	// TLSConfig specifies tls.Config for secure serving
	TlSConfig *tls.Config

	// The interval on which to register
	RegisterInterval time.Duration

	MaxConn int

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

func NewOptions(opts ...Option) Options {
	options := Options{
		Address:          DefaultAddress,
		TlSConfig:        nil,
		RegisterInterval: DefaultRegisterInterval,
		Context:          context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

type Option func(*Options)

func Address(addr string) Option {
	return func(o *Options) {
		o.Address = addr
	}
}

func TLS(tlsCfg *tls.Config) Option {
	return func(o *Options) {
		o.TlSConfig = tlsCfg
	}
}

// MaxConn specifies maximum number of max simultaneous connections to server
func MaxConn(n int) Option {
	return func(o *Options) {
		o.MaxConn = n
	}
}

// Context specifies a context for the rest.
// Can be used to signal shutdown of the rest
// Can be used for extra option values.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

type HandlerKind string

const (
	Resource HandlerKind = "resource"
	Http     HandlerKind = "http"
	RPC      HandlerKind = "rpc"
)

type HandlerOptions struct {
	Kind    HandlerKind
	Target  *ResourceTarget
	Paths   map[string]*pb.OpenAPIPath
	Models  map[string]*pb.Model
	Schemes *pb.SecuritySchemes
}

type ResourceTarget struct {
	Gvk    schema.GroupVersionKind
	Target reflect.Type
}

type HandlerOption func(*HandlerOptions)

func WithKind(kind HandlerKind) HandlerOption {
	return func(o *HandlerOptions) {
		o.Kind = kind
	}
}

func WithPath(url string, path *pb.OpenAPIPath) HandlerOption {
	return func(o *HandlerOptions) {
		if o.Paths == nil {
			o.Paths = map[string]*pb.OpenAPIPath{}
		}
		o.Paths[url] = path
	}
}

func WithModel(path string, model *pb.Model) HandlerOption {
	return func(o *HandlerOptions) {
		if o.Models == nil {
			o.Models = map[string]*pb.Model{}
		}
		o.Models[path] = model
	}
}

func WithSchemes(schemes *pb.SecuritySchemes) HandlerOption {
	return func(o *HandlerOptions) {
		o.Schemes = schemes
	}
}

func WithTarget(gvk schema.GroupVersionKind, target runtime.Object) HandlerOption {
	return func(o *HandlerOptions) {
		o.Target = &ResourceTarget{
			Gvk:    gvk,
			Target: reflect.TypeOf(target),
		}
	}
}

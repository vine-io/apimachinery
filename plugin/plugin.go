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

package plugin

import (
	"errors"
	"fmt"
	"sync"

	"github.com/vine-io/vine/core/server"
)

var (
	// ErrNoType is returned when no type is specified
	ErrNoType = errors.New("plugin: no type")
	// ErrNoPluginID is returned when no id is specified
	ErrNoPluginID = errors.New("plugin: no id")
	// ErrIDRegistered is returned when a duplicate id is already registered
	ErrIDRegistered = errors.New("plugin: id already registered")
	// ErrSkipPlugin is used when a plugin is not initialized and should not be loaded.
	// this allows the plugin loader differentiate between a plugin which is configured
	// not to load and one that fails to load.
	ErrSkipPlugin = errors.New("skip plugin")

	// ErrAlreadyExists is returned when a type already exists
	ErrAlreadyExists = errors.New("already exists")
	// ErrNotFound is returned when a type not be founded
	ErrNotFound = errors.New("not found")

	// ErrInvalidRequires will be thrown if the requirements for a plugin are
	// defined in an invalid manner.
	ErrInvalidRequires = errors.New("invalid requires")
)

// IsSkipPlugin returns true if the error is skipping the plugin
func IsSkipPlugin(err error) bool {
	return errors.Is(err, ErrSkipPlugin)
}

// Type is the type of the plugin
type Type string

func (t Type) String() string { return string(t) }

const (
	// InternalPlugin implements an internal plugin to vine
	InternalPlugin Type = "io.vine.internal.v1"
	// ServicePlugin implements an internal service
	ServicePlugin Type = "io.vine.service.v1"
	// VinePlugin implements a vine service
	VinePlugin Type = "io.vine.vine.v1"
	// GRPCPlugin implements a grpc service
	GRPCPlugin Type = "io.vine.grpc.v1"
)

// Registration contains information for registering a plugin
type Registration struct {
	// Type of the plugin
	Type Type
	// ID of the plugin
	ID string
	// Config specific to the plugin
	Config interface{}
	// Requires is a list of plugins that the registered plugin requires to be available
	Requires []Type

	// InitFn is called when initializing a plugin. The registration and
	// context are passed in. The init function may modify the registration to
	// and exports, capabilities and platform support declarations
	InitFn func(*InitContext) (interface{}, error)
	// Disable the plugin from loading
	Disable bool
}

// Init the registered plugin
func (r *Registration) Init(ic *InitContext) *Plugin {
	p, err := r.InitFn(ic)
	return &Plugin{
		Registration: r,
		Config:       ic.Config,
		Meta:         ic.Meta,
		instance:     p,
		err:          err,
	}
}

// URI returns the full plugin URI
func (r *Registration) URI() string {
	return fmt.Sprintf("%s.%s", r.Type, r.ID)
}

// Service allows Vine services to be registered with the underlying server
type Service interface {
	Register(server.Server) error
}

var register = struct {
	sync.RWMutex
	r []*Registration
}{}

// Load loads all plugins at the provided path into vine service
func Load(path string) (err error) {
	defer func() {
		if v := recover(); v != nil {
			rerr, ok := v.(error)
			if !ok {
				rerr = fmt.Errorf("%s", v)
			}
			err = rerr
		}
	}()
	return loadPlugins(path)
}

// Register allows plugins to register
func Register(r *Registration) {
	register.Lock()
	defer register.Unlock()

	if r.Type == "" {
		panic(ErrNoType)
	}
	if r.ID == "" {
		panic(ErrNoPluginID)
	}
	if err := checkUnique(r); err != nil {
		panic(err)
	}

	var last bool
	for _, requires := range r.Requires {
		if requires == "*" {
			last = true
		}
	}
	if last && len(r.Requires) != 1 {
		panic(ErrInvalidRequires)
	}

	register.r = append(register.r, r)
}

func checkUnique(r *Registration) error {
	for _, registered := range register.r {
		if r.URI() == registered.URI() {
			return fmt.Errorf("%w: %v", ErrIDRegistered, r.URI())
		}
	}
	return nil
}

// DisableFilter filters out disable plugins
type DisableFilter func(r *Registration) bool

// Graph returns an ordered list of registered plugins for initialization.
// Plugins in disableList specified be id will be disabled.
func Graph(filter DisableFilter) (ordered []*Registration) {
	register.RLock()
	defer register.RUnlock()

	for _, r := range register.r {
		if filter(r) {
			r.Disable = true
		}
	}

	added := map[*Registration]bool{}
	for _, r := range register.r {
		if r.Disable {
			continue
		}
		children(r, added, &ordered)
		if !added[r] {
			ordered = append(ordered, r)
			added[r] = true
		}
	}
	return ordered
}

func children(reg *Registration, added map[*Registration]bool, ordered *[]*Registration) {
	for _, t := range reg.Requires {
		for _, r := range register.r {
			if !r.Disable &&
				r.URI() != reg.URI() &&
				(t == "*" || reg.URI() == t) {
				children(r, added, ordered)
				if !added[r] {
					*ordered = append(*ordered, r)
					added[r] = true
				}
			}
		}
	}
}

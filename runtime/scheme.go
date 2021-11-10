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

package runtime

import (
	"fmt"
	"reflect"

	"github.com/vine-io/apimachinery/schema"
)

var DefaultScheme Scheme = NewScheme()

type SimpleScheme struct {
	gvkToTypes map[schema.GroupVersionKind]reflect.Type

	typesToGvk map[reflect.Type]schema.GroupVersionKind

	defaultFuncs map[reflect.Type]DefaultFunc

	observedVersions []schema.GroupVersion
}

// New creates a new Object
func (s *SimpleScheme) New(gvk schema.GroupVersionKind) (Object, error) {
	rv, exists := s.gvkToTypes[gvk]
	if !exists {
		return nil, ErrUnknownGVK
	}

	out := reflect.New(rv).Interface().(Object)

	return out, nil
}

func (s *SimpleScheme) IsExists(gvk schema.GroupVersionKind) bool {
	_, ok := s.gvkToTypes[gvk]
	return ok
}

// AddKnownTypes add Object to Scheme
func (s *SimpleScheme) AddKnownTypes(gv schema.GroupVersion, types ...Object) error {
	s.addObservedVersion(gv)
	for _, v := range types {
		rt := reflect.TypeOf(v)
		if rt.Kind() != reflect.Ptr {
			return ErrIsNotPointer
		}
		gvk := gv.WithKind(rt.Name())
		s.gvkToTypes[gvk] = rt
		s.typesToGvk[rt] = gvk
	}

	return nil
}

func (s *SimpleScheme) Default(src Object) error {
	rt := reflect.TypeOf(src)
	fn, exists := s.defaultFuncs[rt]
	if !exists {
		return fmt.Errorf("%w: not found %v", ErrUnknownType, src.GetObjectKind().GroupVersionKind())
	}
	fn(src)
	return nil
}

func (s *SimpleScheme) AddTypeDefaultingFunc(srcType Object, fn DefaultFunc) {
	s.defaultFuncs[reflect.TypeOf(srcType)] = fn
}

func (s *SimpleScheme) addObservedVersion(gv schema.GroupVersion) {
	if gv.Version == "" {
		return
	}

	for _, observedVersion := range s.observedVersions {
		if observedVersion == gv {
			return
		}
	}

	s.observedVersions = append(s.observedVersions, gv)
}

func NewScheme() *SimpleScheme {
	return &SimpleScheme{
		gvkToTypes:       map[schema.GroupVersionKind]reflect.Type{},
		typesToGvk:       map[reflect.Type]schema.GroupVersionKind{},
		defaultFuncs:     map[reflect.Type]DefaultFunc{},
		observedVersions: []schema.GroupVersion{},
	}
}

// NewObject calls DefaultScheme.New()
func NewObject(gvk schema.GroupVersionKind) (Object, error) {
	return DefaultScheme.New(gvk)
}

// IsExists calls DefaultScheme.IsExists()
func IsExists(gvk schema.GroupVersionKind) bool {
	return DefaultScheme.IsExists(gvk)
}

// AddKnownTypes calls DefaultScheme.AddKnownTypes()
func AddKnownTypes(gv schema.GroupVersion, types ...Object) error {
	return DefaultScheme.AddKnownTypes(gv, types...)
}

// DefaultObject calls DefaultScheme.Default()
func DefaultObject(src Object) error {
	return DefaultScheme.Default(src)
}

// AddTypeDefaultingFunc calls DefaultScheme.AddTypeDefaultingFunc()
func AddTypeDefaultingFunc(srcType Object, fn DefaultFunc) {
	DefaultScheme.AddTypeDefaultingFunc(srcType, fn)
}

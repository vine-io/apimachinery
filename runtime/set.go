package runtime

import (
	"sync"

	"github.com/vine-io/apimachinery/schema"
)

var DefaultSet ObjectSet = NewSet()

type SimpleObjectSet struct {
	sync.RWMutex

	store map[string]Object
}

// NewObj creates a new Object, trigger OnCreate function
func (s *SimpleObjectSet) NewObj(gvk string, fn Creater) (Object, bool) {
	s.RLock()
	defer s.RUnlock()
	out, ok := s.store[gvk]
	if !ok {
		return nil, false
	}

	out = out.DeepCopyObject()
	if fn == nil {
		return out, true
	}

	return fn(out), true
}

// NewObjWithGVK creates a new Object, trigger OnCreate function
func (s *SimpleObjectSet) NewObjWithGVK(gvk *schema.GroupVersionKind, fn Creater) (Object, bool) {
	s.RLock()
	defer s.RUnlock()
	out, ok := s.store[gvk.String()]
	if !ok {
		return nil, false
	}

	out = out.DeepCopyObject()
	if fn == nil {
		return out, true
	}

	return fn(out), true
}

func (s *SimpleObjectSet) IsExists(gvk *schema.GroupVersionKind) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.store[gvk.String()]
	return ok
}

// Get gets Object by schema.GroupVersionKind
func (s *SimpleObjectSet) Get(gvk *schema.GroupVersionKind) (Object, bool) {
	s.RLock()
	defer s.RUnlock()
	out, ok := s.store[gvk.String()]
	if !ok {
		return nil, false
	}
	return out.DeepCopyObject(), true
}

// AddObj push entities to Set
func (s *SimpleObjectSet) AddObj(v ...Object) {
	s.Lock()
	for _, in := range v {
		gvk := in.GetObjectKind().GroupVersionKind()
		s.store[gvk.String()] = in
	}
	s.Unlock()
}

// NewObj creates a new Object, trigger OnCreate function
func NewObj(gvk string, fn Creater) (Object, bool) {
	return DefaultSet.NewObj(gvk, fn)
}

// NewObjWithGVK creates a new Object, trigger OnCreate function
func NewObjWithGVK(gvk *schema.GroupVersionKind, fn Creater) (Object, bool) {
	return DefaultSet.NewObjWithGVK(gvk, fn)
}

func AddObj(vs ...Object) {
	DefaultSet.AddObj(vs...)
}

func NewSet() *SimpleObjectSet {
	return &SimpleObjectSet{
		store: map[string]Object{},
	}
}

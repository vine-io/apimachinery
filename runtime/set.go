package runtime

import (
	"sync"

	"github.com/vine-io/apimachinery/schema"
)

var gs = NewObjectSet()

type Creater func(in Object) Object

var DefaultCreater Creater = func(in Object) Object {
	return in
}

type ObjectSet struct {
	sync.RWMutex

	store map[string]Object
}

// NewObj creates a new Object, trigger OnCreate function
func (s *ObjectSet) NewObj(gvk string, fn Creater) (Object, bool) {
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
func (s *ObjectSet) NewObjWithGVK(gvk *schema.GroupVersionKind, fn Creater) (Object, bool) {
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

func (s *ObjectSet) IsExists(gvk *schema.GroupVersionKind) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.store[gvk.String()]
	return ok
}

// Get gets Object by schema.GroupVersionKind
func (s *ObjectSet) Get(gvk *schema.GroupVersionKind) (Object, bool) {
	s.RLock()
	defer s.RUnlock()
	out, ok := s.store[gvk.String()]
	if !ok {
		return nil, false
	}
	return out.DeepCopyObject(), true
}

// AddObj push entities to Set
func (s *ObjectSet) AddObj(v ...Object) {
	s.Lock()
	for _, in := range v {
		gvk := in.GetObjectKind().GroupVersionKind()
		s.store[gvk.String()] = in
	}
	s.Unlock()
}

// NewObj creates a new Object, trigger OnCreate function
func NewObj(gvk string, fn Creater) (Object, bool) {
	return gs.NewObj(gvk, fn)
}

// NewObjWithGVK creates a new Object, trigger OnCreate function
func NewObjWithGVK(gvk *schema.GroupVersionKind, fn Creater) (Object, bool) {
	return gs.NewObjWithGVK(gvk, fn)
}

func AddObj(vs ...Object) {
	gs.AddObj(vs...)
}

func NewObjectSet() *ObjectSet {
	return &ObjectSet{
		store: map[string]Object{},
	}
}

package store

import (
	"sync"

	"github.com/vine-io/apimachinery/runtime"
	"github.com/vine-io/apimachinery/schema"
)

var DefaultSet RepoSet = NewSimpleRepoSet()

type SimpleRepoSet struct {
	sync.RWMutex

	sets map[string]RepoCreator
}

func (s *SimpleRepoSet) RegisterRepo(g *schema.GroupVersionKind, fn RepoCreator) {
	s.Lock()
	s.sets[g.String()] = fn
	s.Unlock()
}

func (s *SimpleRepoSet) NewRepo(in runtime.Object) (Repo, bool) {
	s.RLock()
	defer s.RUnlock()
	gvk := in.GetObjectKind().GroupVersionKind()
	fn, ok := s.sets[gvk.String()]
	if !ok {
		return nil, false
	}
	return fn(in), true
}

func (s *SimpleRepoSet) IsExits(gvk *schema.GroupVersionKind) bool {
	s.RLock()
	_, ok := s.sets[gvk.String()]
	s.RUnlock()
	return ok
}

func NewSimpleRepoSet() *SimpleRepoSet {
	return &SimpleRepoSet{sets: map[string]RepoCreator{}}
}

func RegisterRepo(g *schema.GroupVersionKind, fn RepoCreator) {
	DefaultSet.RegisterRepo(g, fn)
}

func NewRepo(in runtime.Object) (Repo, bool) {
	return DefaultSet.NewRepo(in)
}

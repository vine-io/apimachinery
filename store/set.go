package store

import (
	"sync"

	"github.com/vine-io/apimachinery/runtime"
	"github.com/vine-io/apimachinery/schema"
)

var pset = NewRepoSet()

type RepoCreator func(object runtime.Object) Repo

type RepoSet struct {
	sync.RWMutex

	sets map[string]RepoCreator
}

func (s *RepoSet) RegisterRepo(g *schema.GroupVersionKind, fn RepoCreator) {
	s.Lock()
	s.sets[g.String()] = fn
	s.Unlock()
}

func (s *RepoSet) NewRepo(in runtime.Object) (Repo, bool) {
	s.RLock()
	defer s.RUnlock()
	gvk := in.GetObjectKind().GroupVersionKind()
	fn, ok := s.sets[gvk.String()]
	if !ok {
		return nil, false
	}
	return fn(in), true
}

func (s *RepoSet) IsExits(gvk *schema.GroupVersionKind) bool {
	s.RLock()
	_, ok := s.sets[gvk.String()]
	s.RUnlock()
	return ok
}

func NewRepoSet() *RepoSet {
	return &RepoSet{sets: map[string]RepoCreator{}}
}

func RegisterRepo(g *schema.GroupVersionKind, fn RepoCreator) {
	pset.RegisterRepo(g, fn)
}

func NewRepo(in runtime.Object) (Repo, bool) {
	return pset.NewRepo(in)
}

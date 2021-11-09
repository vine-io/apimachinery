package storage

import (
	"fmt"
	"reflect"

	"github.com/vine-io/apimachinery/runtime"
	"github.com/vine-io/apimachinery/schema"
)

var (
	ErrInvalidObject       = fmt.Errorf("specified object invalid")
	ErrStorageIsNotPointer = fmt.Errorf("storage is not a pointer")
	ErrStorageNotExists    = fmt.Errorf("storage not exists")
	ErrStorageAutoMigrate  = fmt.Errorf("auto migrate storage")
)

var DefaultFactory Factory = NewStorageFactory()

type simpleStorageFactory struct {
	gvkToType map[schema.GroupVersionKind]reflect.Type
}

func (s *simpleStorageFactory) AddKnownStorage(gvk schema.GroupVersionKind, storage Storage) error {
	rt := reflect.TypeOf(storage)
	if rt.Kind() != reflect.Ptr {
		return ErrStorageIsNotPointer
	}

	s.gvkToType[gvk] = rt

	if err := storage.AutoMigrate(); err != nil {
		return fmt.Errorf("%w: %v", ErrStorageAutoMigrate, err)
	}

	return nil
}

func (s *simpleStorageFactory) NewStorage(in runtime.Object) (Storage, error) {
	gvk := in.GetObjectKind().GroupVersionKind()
	rt, exists := s.gvkToType[gvk]
	if !exists {
		return nil, fmt.Errorf("%w: object's gvk is %s", ErrStorageNotExists, in.GetObjectKind().GroupVersionKind())
	}

	storage := reflect.New(rt).Interface().(Storage)

	return storage, nil
}

func (s *simpleStorageFactory) IsExists(gvk schema.GroupVersionKind) bool {
	_, ok := s.gvkToType[gvk]
	return ok
}

func NewStorageFactory() Factory {
	return &simpleStorageFactory{
		gvkToType: map[schema.GroupVersionKind]reflect.Type{},
	}
}

func AddKnownStorage(gvk schema.GroupVersionKind, storage Storage) error {
	return DefaultFactory.AddKnownStorage(gvk, storage)
}

func NewStorage(in runtime.Object) (Storage, error) {
	return DefaultFactory.NewStorage(in)
}

func IsExists(gvk schema.GroupVersionKind) bool {
	return DefaultFactory.IsExists(gvk)
}

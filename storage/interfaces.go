package storage

import (
	"context"

	"github.com/vine-io/apimachinery/runtime"
	"github.com/vine-io/apimachinery/schema"
	"github.com/vine-io/vine/lib/dao"
	"github.com/vine-io/vine/lib/dao/clause"
)

type Storage interface {
	AutoMigrate() error
	Load(object runtime.Object) error
	FindPage(ctx context.Context, page, size int32) ([]runtime.Object, int64, error)
	FindAll(ctx context.Context) ([]runtime.Object, error)
	FindPureAll(ctx context.Context) ([]runtime.Object, error)
	Count(ctx context.Context) (total int64, err error)
	FindOne(ctx context.Context) (runtime.Object, error)
	FindPureOne(ctx context.Context) (runtime.Object, error)
	Cond(exprs ...clause.Expression) Storage
	Create(ctx context.Context) (runtime.Object, error)
	BatchUpdates(ctx context.Context) error
	Updates(ctx context.Context) (runtime.Object, error)
	BatchDelete(ctx context.Context, soft bool) error
	Delete(ctx context.Context, soft bool) error
	Tx(ctx context.Context) *dao.DB
}

type Factory interface {
	// AddKnownStorage registers Storage
	AddKnownStorage(gvk schema.GroupVersionKind, storage Storage) error

	// NewStorage get a Storage by runtime.Object
	NewStorage(in runtime.Object) (Storage, error)

	// IsExists checks Storage exists
	IsExists(gvk schema.GroupVersionKind) bool
}

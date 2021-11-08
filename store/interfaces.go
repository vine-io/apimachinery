package store

import (
	"context"

	"github.com/vine-io/apimachinery/runtime"
	"github.com/vine-io/apimachinery/schema"
	"github.com/vine-io/vine/lib/dao"
	"github.com/vine-io/vine/lib/dao/clause"
)

type RepoCreator func(object runtime.Object) Repo

type Repo interface {
	FindPage(ctx context.Context, page, size int32) ([]runtime.Object, int64, error)
	FindAll(ctx context.Context) ([]runtime.Object, error)
	FindPureAll(ctx context.Context) ([]runtime.Object, error)
	Count(ctx context.Context) (total int64, err error)
	FindOne(ctx context.Context) (runtime.Object, error)
	FindPureOne(ctx context.Context) (runtime.Object, error)
	Cond(exprs ...clause.Expression) Repo
	Create(ctx context.Context) (runtime.Object, error)
	BatchUpdates(ctx context.Context) error
	Updates(ctx context.Context) (runtime.Object, error)
	BatchDelete(ctx context.Context, soft bool) error
	Delete(ctx context.Context, soft bool) error
	Tx(ctx context.Context) *dao.DB
}

type RepoSet interface {
	RegisterRepo(g *schema.GroupVersionKind, fn RepoCreator)
	NewRepo(in runtime.Object) (Repo, bool)
	IsExits(gvk *schema.GroupVersionKind) bool
}

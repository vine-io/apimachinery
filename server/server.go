package server

import (
	"context"

	"github.com/vine-io/vine/core/server"
)

// Service the implementation of vine service
type Service interface {
	Register(s server.Server) error
}

// BizImpl the implementation of business service
type BizImpl interface {
	Init(ctx context.Context) error
	Start() error
}

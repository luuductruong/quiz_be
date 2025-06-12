package context

import (
	"context"
	"github.com/quiz_be/services/core/middleware"
	"gorm.io/gorm"
)

type ctxInternal struct {
	context.Context
}

func FromContext(ctx context.Context) Context {
	return &ctxInternal{ctx}
}

type Context interface {
	context.Context
	// from context
	GetDbTx() *gorm.DB
	GetTracerId() string
	Clone() Context
}

func (c *ctxInternal) Clone() Context {
	newCtx := context.Background()
	// Preserve tracer ID
	if tracer := c.GetTracerId(); tracer != "" {
		newCtx = middleware.TracerToContext(newCtx, tracer)
	}
	// Preserve DB transaction
	if db := c.GetDbTx(); db != nil {
		// *gorm.DB internally holds a context.
		// If we directly put the existing db into newCtx, it will carry the old (possibly canceled) context,
		// which can cause "context canceled" errors when used later.
		// Therefore, we create a new copy of *gorm.DB with the new context using db.WithContext(newCtx),
		// and store this copy into newCtx to avoid using the canceled context.
		//newCtx = middleware.DbTxToContext(newCtx, db)
		newCtx = middleware.DbTxToContext(newCtx, db.WithContext(newCtx))
	}
	return &ctxInternal{newCtx}
}

func (c *ctxInternal) GetDbTx() *gorm.DB {
	return middleware.DbTxFromContext(c)
}

func (c *ctxInternal) GetTracerId() string {
	return middleware.TracerFromContext(c)
}

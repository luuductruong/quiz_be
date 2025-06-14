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
	GetPageAndLimit() (int, int)
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

func (c *ctxInternal) GetPageAndLimit() (int, int) {
	return c.GetPage(), c.GetLimit()
}

func (c *ctxInternal) GetPage() int {
	page := 1
	pageCtx := middleware.PageFromCtx(c)
	if pageCtx != nil {
		page = pageCtx.GetPage()
	}
	if page < 1 {
		page = 1
	}
	return page
}

func (c *ctxInternal) GetLimit() int {
	limit := 20
	pageCtx := middleware.PageFromCtx(c)
	if pageCtx != nil {
		limit = pageCtx.GetLimit()
	}
	if limit == -1 {
		limit = 20 // default
	} else if limit == 0 {
		limit = 1
	} else if limit > 100 {
		limit = 100
	}
	return limit
}

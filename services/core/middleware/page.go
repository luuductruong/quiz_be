package middleware

import (
	"context"
	"google.golang.org/grpc"
)

const (
	PageCtxKey = "page"
)

type pageModel struct {
	page  int
	limit int
}

type PageCtx interface {
	GetPage() int
	GetLimit() int
}

type Pagable interface {
	GetPage() int32
	GetLimit() int32
}

func (p *pageModel) GetPage() int {
	return p.page
}

func (p *pageModel) GetLimit() int {
	return p.limit
}

func PageFromCtx(ctx context.Context) PageCtx {
	raw, ok := ctx.Value(PageCtxKey).(*pageModel)
	if !ok {
		return nil
	}
	return raw
}

func RequestPaging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		p := new(pageModel)
		pagable, ok := req.(Pagable)
		if ok {
			p.page = int(pagable.GetPage())
			p.limit = int(pagable.GetLimit())

		}
		ctx = context.WithValue(ctx, PageCtxKey, p)
		return handler(ctx, req)
	}
}

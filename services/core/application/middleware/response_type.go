package middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net/http"
)

const (
	ResponseTypeCtxKey = "ResponseTypeCtxKey"
)

func GrpcResponseType() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ctx, req)
	}
}

func WithResponseType() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// any
			ctx := r.Context()
			responseTypeFromRequest := r.URL.Query().Get("response_type")
			rsType := responseTypeFromRequest

			ctx = metadata.AppendToOutgoingContext(ctx, ResponseTypeCtxKey, rsType)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

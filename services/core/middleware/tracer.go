package middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/quiz_be/services/core/helper"
	"github.com/quiz_be/services/core/infra/trace"
)

const TracerCtxKey = "tracerCtx"

// GrpcTracerId create and add tracerId to context
func GrpcTracerId() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		tracer := md.Get(TracerCtxKey)
		if len(tracer) == 0 {
			traceId := trace.TraceIdFromContext(ctx)
			if traceId == "" || traceId == "00000000000000000000000000000000" {
				tracer = []string{helper.RandString(16)}
			} else {
				tracer = []string{traceId}
			}
		}
		ctx = context.WithValue(ctx, TracerCtxKey, tracer[0])
		ctx = metadata.AppendToOutgoingContext(ctx, TracerCtxKey, tracer[0])
		header := metadata.Pairs(TracerCtxKey, tracer[0])
		err := grpc.SetHeader(ctx, header)
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func TracerFromContext(ctx context.Context) string {
	raw, _ := ctx.Value(TracerCtxKey).(string)
	return raw
}

func TracerToContext(ctx context.Context, tracerId string) context.Context {
	ctx = context.WithValue(ctx, TracerCtxKey, tracerId)
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		md[TracerCtxKey] = []string{tracerId}
		return metadata.NewOutgoingContext(ctx, md)
	}
	return metadata.AppendToOutgoingContext(ctx, TracerCtxKey, tracerId)
}

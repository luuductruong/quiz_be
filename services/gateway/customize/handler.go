package customize

import (
	"context"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/proto"

	baseHttp "github.com/quiz_be/services/core/http"
	"github.com/quiz_be/services/core/middleware"
)

func HttpErrorHandler(ctx context.Context, mux *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		grpclog.Infof("Failed to extract ServerMetadata from context")
	}
	//add expose header
	for k, vs := range md.HeaderMD {
		if h, ok := ExposeHeader(k); ok {
			for _, v := range vs {
				w.Header().Add(h, v)
			}
		}
	}
	baseHttp.WriteError(w, err)
}

func HttpSuccessHandler(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	return baseHttp.ResponseJson(w, resp)
}

func ExposeHeader(key string) (string, bool) {
	key = strings.ToLower(key)
	switch key {
	case strings.ToLower(middleware.TracerCtxKey):
		return "tracer", true
	default:
		return key, false
	}
}

func HeaderAllows(key string) (string, bool) {
	key = strings.ToLower(key)
	switch key {
	case "accept-language":
		return key, true
	case "tracer":
		return key, true
	case "signature":
		return key, true
	case "client-id":
		return key, true
	case "user-id":
		return key, true
	default:
		return key, false
	}
}

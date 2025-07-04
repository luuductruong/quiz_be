package http

import (
	"context"
	"encoding/json"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"net/http"
	"strings"

	appContext "github.com/quiz_be/services/core/context"
	"github.com/quiz_be/services/core/infra/logger"
)

type dataResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
}

type errorResponse struct {
	Success bool          `json:"success"`
	Errors  []interface{} `json:"errors"`
}

var (
	marshaller = &protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}
)

type HandlerFn = func(ctx appContext.Context, req proto.Message, opts ...grpc.CallOption) (proto.Message, error)

type BaseHandler interface {
	HandleReqWithMd(w http.ResponseWriter, req *http.Request, msg proto.Message, handlerFn HandlerFn)

	ResponseJson(w http.ResponseWriter, msg proto.Message, statusCode ...int)
	ResponseError(w http.ResponseWriter, err error, statusCode ...int)
}

func NewBaseHandler() BaseHandler {
	return &baseHandler{logger: logger.Default}
}

type baseHandler struct {
	logger logger.Logger
}

func (handler *baseHandler) HandleReqWithMd(w http.ResponseWriter, req *http.Request, msg proto.Message, handlerFn HandlerFn) {
	var med runtime.ServerMetadata
	ctx := req.Context()
	for k, v := range req.Header {
		if strings.ToLower(k) == "authorization" {
			ctx = context.WithValue(ctx, strings.ToLower(k), v[0])
			ctx = metadata.AppendToOutgoingContext(ctx, strings.ToLower(k), v[0])
		}
		if strings.ToLower(k) == "client-id" {
			ctx = context.WithValue(ctx, strings.ToLower(k), v[0])
			ctx = metadata.AppendToOutgoingContext(ctx, strings.ToLower(k), v[0])
		}
		if strings.ToLower(k) == "accept-language" {
			ctx = context.WithValue(ctx, strings.ToLower(k), v[0])
			ctx = metadata.AppendToOutgoingContext(ctx, strings.ToLower(k), v[0])
		}
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "from_gateway", "true")
	ctx = metadata.AppendToOutgoingContext(ctx, "user-agent", "123")
	appCtx := appContext.FromContext(ctx)
	res, err := handlerFn(appCtx, msg, grpc.Header(&med.HeaderMD), grpc.Trailer(&med.TrailerMD))
	setDefaultHeaderValues(w, med)
	if err != nil {
		handler.logger.ErrorCtx(appCtx, err)
		handler.ResponseError(w, err)
		return
	}
	handler.ResponseJson(w, res)
}

func (handler *baseHandler) ResponseError(w http.ResponseWriter, err error, statusCode ...int) {
	WriteError(w, err, statusCode...)
}

func (handler *baseHandler) ResponseJson(w http.ResponseWriter, data proto.Message, statusCode ...int) {
	var httpCode = http.StatusOK
	if len(statusCode) > 0 {
		httpCode = statusCode[0]
	}

	dataPayload, err := marshaller.Marshal(data)
	if err != nil {
		handler.logger.Error(err)
		return
	}

	body, err := json.Marshal(dataResponse{Success: true, Data: []byte(dataPayload)})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(body)
}

func ResponseJson(w http.ResponseWriter, data proto.Message, statusCode ...int) error {
	var httpCode = http.StatusOK
	if len(statusCode) > 0 {
		httpCode = statusCode[0]
	}
	dataPayload, err := marshaller.Marshal(data)
	if err != nil {
		return err
	}

	body, err := json.Marshal(dataResponse{Success: true, Data: []byte(dataPayload)})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(body)
	return nil
}

func WriteError(w http.ResponseWriter, err error, statusCode ...int) {
	errStatus, ok := status.FromError(err)
	if !ok {
		errStatus = status.New(codes.Unknown, err.Error())
	}

	var httpCode = http.StatusBadRequest
	if len(statusCode) > 0 {
		httpCode = statusCode[0]
	} else {
		httpCode = runtime.HTTPStatusFromCode(errStatus.Code())
	}

	body, err := json.Marshal(errorResponse{Success: false, Errors: errStatus.Details()})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(body)
}

func setDefaultHeaderValues(w http.ResponseWriter, md runtime.ServerMetadata) {
	for k, v := range md.HeaderMD {
		if strings.EqualFold(k, "tracer") {
			w.Header().Set(k, v[0])
			break
		}
	}
}

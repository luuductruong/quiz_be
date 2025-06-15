package cmd

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/quiz_be/services/gateway/application"
	"github.com/quiz_be/services/gateway/application/http/ws"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"time"

	coreMdw "github.com/quiz_be/services/core/application/middleware"
	quizService "github.com/quiz_be/services/core/application/quiz/service"
	"github.com/quiz_be/services/core/infra/config"
	grpcHelper "github.com/quiz_be/services/core/infra/grpc"
	"github.com/quiz_be/services/core/infra/logger"
	gatewayMdw "github.com/quiz_be/services/gateway/middleware"
)

var (
	appConfig *Config

	// client connect
	quizConn *grpc.ClientConn
	// end client connect
	hub *ws.WsHub
)

func Run() {
	var err error
	hub = ws.NewWsHub()
	config.LoadConfig(&appConfig)
	log := logger.NewLogger(appConfig.Logger)
	logger.SetDefault(log)

	// client connect
	quizConn, err = grpcHelper.NewClient(appConfig.Quiz)
	if err != nil {
		logger.Default.Panic("can't connect to quiz: ", err)
	}
	// end client connect

	// make context with cancel
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// setup timeout for request
	runtime.DefaultContextTimeout = 30 * time.Second // for example

	gMux := runtime.NewServeMux(
		// https://grpc-ecosystem.github.io/grpc-gateway/docs/development/grpc-gateway_v2_migration_guide/
		// use UseProtoNames = true for JSON snake_case (default is camelCase)
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}))

	// register client connect
	err = quizService.RegisterQuizServiceHandler(ctx, gMux, quizConn)
	if err != nil {
		logger.Default.Panic("can't register quiz: ", err)
	}
	// end register client connect
	httpHandler := application.NewHttpHandler(quizConn, hub)

	gMux.HandlePath("POST", "/v1/quiz/submit-answer", httpHandler.Quiz().SubmitAnswer)
	gMux.HandlePath("GET", "/v1/quiz/{quiz_id}/detail", httpHandler.Quiz().GetQuizDetail)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", hub.HandleWebSocket)
	mux.Handle("/", gatewayMdw.WithCORS(gMux))
	combineMd := gatewayMdw.ChainCombine(coreMdw.WithResponseType())
	server := http.Server{
		Handler: combineMd(mux),
	}
	listen, err := net.Listen("tcp", appConfig.Http.Address())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Gateway start listening on " + appConfig.Http.Address())
	log.Fatal(server.Serve(listen))
}

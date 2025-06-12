package cmd

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/quiz_be/services/core/infra/logger"
	"net"
	"net/http"
	"time"

	coreMdw "github.com/quiz_be/services/core/application/middleware"
	"github.com/quiz_be/services/core/infra/config"
	gatewayMdw "github.com/quiz_be/services/gateway/middleware"
)

var (
	appConfig *Config

	// client connect
)

func Run() {
	var err error
	config.LoadConfig(&appConfig)
	log := logger.NewLogger(appConfig.Logger)
	logger.SetDefault(log)

	// client connect

	// end client connect

	// make context with cancel
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// setup timeout for request
	runtime.DefaultContextTimeout = 30 * time.Second // for example

	gMux := runtime.NewServeMux()

	// register client connect

	// end register client connect

	mux := http.NewServeMux()
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

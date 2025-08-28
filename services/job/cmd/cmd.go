package cmd

import (
	"fmt"
	"github.com/quiz_be/services/core/domain/job"
	"github.com/quiz_be/services/core/infra/factory"
	"github.com/quiz_be/services/job/external/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"

	appService "github.com/quiz_be/services/core/application/job/service"
	"github.com/quiz_be/services/core/application/middleware"
	"github.com/quiz_be/services/core/i18n"
	"github.com/quiz_be/services/core/infra/config"
	"github.com/quiz_be/services/core/infra/db"
	"github.com/quiz_be/services/core/infra/logger"
	"github.com/quiz_be/services/core/infra/pubsub"
	handler "github.com/quiz_be/services/job/application/grpchandler"
	"github.com/quiz_be/services/job/domain"
)

var (
	jobDomain     job.Service
	sql           db.SQL
	appConfig     *Config
	grpcHandler   appService.JobServiceServer
	appSubscriber pubsub.AppSubscriber
	psClient      pubsub.PubSub
	messageBus    pubsub.MessageBus
)

func Run() {
	fmt.Println("hello, my name is Job")
	var err error
	//appConfig = &Config{}
	err = config.LoadConfig(&appConfig)
	if err != nil {
		logger.Default.Panic("Error loading config: ", err)
	}
	log := logger.NewLogger(appConfig.Logger)
	logger.SetDefault(log)
	sql, err = db.NewSQL(appConfig.DataBase)
	if err != nil {
		logger.Default.Panic("Error connecting to database: ", err)
	}
	i18n.Init(appConfig.I18n)
	messageBus = pubsub.NewBatchStagedMessageBus(appConfig.PubSub.Topic)
	psClient, err = factory.NewPubSub(log, appConfig.PubSub)
	if err != nil {
		logger.Default.Panic("can't connect to pubsub: ", err)
	}

	jobDomain = domain.NewDomain(&domain.JobDomainParam{
		Publisher: psClient.Publisher(),
		JobRepo:   repository.NewJobRepo(),
	})
	grpcHandler = handler.NewHandler(jobDomain)
	appSubscriber = pubsub.NewAppSubscriber(psClient.Subscriber(), appConfig.PubSub.Subscription, nil)
	//subsHandler := subscriber.NewHandler(&subscriber.SubsHandlerParam{
	//	Service: jobDomain,
	//	DB:      sql.GetDB(),
	//	Logger:  log,
	//})
	//appSubscriber.RegisterEventSubscriber(subsHandler.RouteSetup())
	//err = appSubscriber.StartReceiving()
	//if err != nil {
	//	logger.Default.Panic("can't start receiving message: ", err)
	//}
	//defer appSubscriber.StopReceiving()
	grpcServe()
}

func grpcServe() {
	// Start gRPC server in goroutine
	lis, err := net.Listen("tcp", appConfig.Grpc.Address())
	if err != nil {
		logger.Default.Panic("failed to listen: %v", err)
	}
	fmt.Println("grpc server listening on :", lis.Addr().String())
	serve := middleware.NewServer(grpc.UnaryInterceptor(middleware.GrpcChainUnaryServer(sql)))
	defer serve.GracefulStop()

	appService.RegisterJobServiceServer(serve, grpcHandler)
	reflection.Register(serve)
	logger.Default.Fatal(serve.Serve(lis))
}

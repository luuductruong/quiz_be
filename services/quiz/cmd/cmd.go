package cmd

import (
	"fmt"
	jobClient "github.com/quiz_be/services/core/client/job"
	"github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/infra/factory"
	grpcHelper "github.com/quiz_be/services/core/infra/grpc"
	"github.com/quiz_be/services/quiz/application/subscriber"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"

	"github.com/quiz_be/services/core/application/middleware"
	appService "github.com/quiz_be/services/core/application/quiz/service"
	"github.com/quiz_be/services/core/i18n"
	"github.com/quiz_be/services/core/infra/config"
	"github.com/quiz_be/services/core/infra/db"
	"github.com/quiz_be/services/core/infra/logger"
	"github.com/quiz_be/services/core/infra/pubsub"
	handler "github.com/quiz_be/services/quiz/application/grpchandler"
	"github.com/quiz_be/services/quiz/domain"
	repo "github.com/quiz_be/services/quiz/external/repository"
)

var (
	quizDomain    quiz.Service
	sql           db.SQL
	appConfig     *Config
	grpcHandler   appService.QuizServiceServer
	appSubscriber pubsub.AppSubscriber
	psClient      pubsub.PubSub
	messageBus    pubsub.MessageBus
)

func Run() {
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
	//client connect
	jobConn, err := grpcHelper.NewClient(appConfig.Job)
	if err != nil {
		logger.Default.Panic("can't connect to job service: ", err)
	}
	jobService := jobClient.NewService(jobConn)

	quizDomain = domain.NewDomain(&domain.QuizDomainParam{
		QuizRepo:         repo.NewQuizRepo(),
		QuestionRepo:     repo.NewQuestionRepo(),
		QuizQuestionRepo: repo.NewQuizQuestionRepo(),
		UserRepo:         repo.NewUserRepo(),
		UserAnswerRepo:   repo.NewUserAnswerRepo(),
		ScoreRepo:        repo.NewScoreRepo(),
		Publisher:        psClient.Publisher(),
		JobClient:        jobService,
	})
	grpcHandler = handler.NewHandler(quizDomain)
	appSubscriber = pubsub.NewAppSubscriber(psClient.Subscriber(), appConfig.PubSub.Subscription, nil)
	subsHandler := subscriber.NewHandler(&subscriber.SubsHandlerParam{
		Service: quizDomain,
		DB:      sql.GetDB(),
		Logger:  log,
	})
	appSubscriber.RegisterEventSubscriber(subsHandler.RouteSetup())
	err = appSubscriber.StartReceiving()
	if err != nil {
		logger.Default.Panic("can't start receiving message: ", err)
	}
	defer appSubscriber.StopReceiving()
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

	appService.RegisterQuizServiceServer(serve, grpcHandler)
	reflection.Register(serve)
	logger.Default.Fatal(serve.Serve(lis))
}

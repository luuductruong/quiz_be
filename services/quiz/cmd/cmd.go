package cmd

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"

	"github.com/quiz_be/services/core/application/middleware"
	appService "github.com/quiz_be/services/core/application/quiz/service"
	"github.com/quiz_be/services/core/infra/config"
	"github.com/quiz_be/services/core/infra/db"
	"github.com/quiz_be/services/core/infra/logger"
	handler "github.com/quiz_be/services/quiz/application/grpchandler"
	"github.com/quiz_be/services/quiz/domain"
	repo "github.com/quiz_be/services/quiz/external/repository"
)

var (
	sql         db.SQL
	appConfig   *Config
	grpcHandler appService.QuizServiceServer
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

	quizDomain := domain.NewDomain(&domain.QuizDomainParam{
		QuizRepo:         repo.NewQuizRepo(),
		QuestionRepo:     repo.NewQuestionRepo(),
		QuizQuestionRepo: repo.NewQuizQuestionRepo(),
	})
	grpcHandler = handler.NewHandler(quizDomain)

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

package job

import (
	"google.golang.org/grpc"

	"github.com/quiz_be/services/core/application/job/dto"
	jobService "github.com/quiz_be/services/core/application/job/service"
	"github.com/quiz_be/services/core/context"
	"github.com/quiz_be/services/core/infra/logger"
)

type Service interface {
	PushJob(ctx context.Context, name string, topics []string, data []byte) error
}

func NewService(conn *grpc.ClientConn) Service {
	return &service{
		service: jobService.NewJobServiceClient(conn),
		logger:  logger.Default,
	}
}

type service struct {
	service jobService.JobServiceClient
	logger  logger.Logger
}

func (s *service) PushJob(ctx context.Context, name string, topics []string, data []byte) error {
	s.logger.DebugCtx(ctx, "PushJob")
	_, err := s.service.PushJob(ctx, &dto.PushJobReq{
		Name:   name,
		Topics: topics,
		Data:   data,
	})
	if err != nil {
		s.logger.ErrorCtx(ctx, err, "PushJob failed")
		return err
	}
	return nil
}

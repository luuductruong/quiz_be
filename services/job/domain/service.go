package domain

import (
	"github.com/google/uuid"
	"github.com/quiz_be/services/core/context"
	"github.com/quiz_be/services/core/domain/job"
	"github.com/quiz_be/services/core/infra/logger"
	"github.com/quiz_be/services/core/infra/pubsub"
	"time"
)

type JobDomainParam struct {
	Publisher pubsub.Publisher
	JobRepo   job.JobRepo
}

type domain struct {
	logger    logger.Logger
	publisher pubsub.Publisher
	jobRepo   job.JobRepo
}

func (d *domain) PushJob(ctx context.Context, topic string, data []byte) error {
	d.logger.DebugCtx(ctx, "PushJob")
	d.logger.DebugCtx(ctx, "topic: ", topic)
	d.logger.DebugCtx(ctx, "data: ", string(data))
	err := d.jobRepo.Upsert(ctx, &job.Job{
		ID:        uuid.NewString(),
		Name:      topic,
		Data:      data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "upsert job failed")
		return err
	}
	return nil
}

func (d *domain) Test(ctx context.Context, name string) (*job.Test, error) {
	d.logger.DebugCtx(ctx, "Test")
	return &job.Test{Name: name}, nil
}

func NewDomain(param *JobDomainParam) job.Service {
	return &domain{
		logger:    logger.Default,
		publisher: param.Publisher,
		jobRepo:   param.JobRepo,
	}
}

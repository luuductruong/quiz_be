package domain

import (
	"github.com/quiz_be/services/core/context"
	"github.com/quiz_be/services/core/domain/job"
	"github.com/quiz_be/services/core/helper"
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

func (d *domain) PushJob(ctx context.Context, name string, topics []string, data []byte) error {
	d.logger.DebugCtx(ctx, "PushJob")
	d.logger.DebugCtx(ctx, "topics: ", topics)
	d.logger.DebugCtx(ctx, "data: ", string(data))
	var err error
	defer func() {
		err = nil
	}()
	tracerId := ctx.GetTracerId()
	for _, topic := range topics {
		err = d.jobRepo.Upsert(ctx, &job.Job{
			ID:         helper.NewStringUUIDV7(),
			Name:       name,
			Data:       data,
			RetryCount: 0,
			LastError:  "",
			TraceID:    tracerId,
			Exchange:   topic,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})
		if err != nil {
			d.logger.ErrorCtx(ctx, err, "upsert job failed for name: ", name, " topic: ", topic, " data: ", string(data))
			// ignore return error
		}
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

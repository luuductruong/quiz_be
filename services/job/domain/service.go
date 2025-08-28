package domain

import (
	"github.com/quiz_be/services/core/context"
	"github.com/quiz_be/services/core/domain/job"
	"github.com/quiz_be/services/core/infra/logger"
	"github.com/quiz_be/services/core/infra/pubsub"
)

type JobDomainParam struct {
	Publisher pubsub.Publisher
}

type domain struct {
	logger    logger.Logger
	publisher pubsub.Publisher
}

func (d *domain) Test(ctx context.Context, name string) (*job.Test, error) {
	d.logger.DebugCtx(ctx, "Test")
	return &job.Test{Name: name}, nil
}

func NewDomain(param *JobDomainParam) job.Service {
	return &domain{
		logger:    logger.Default,
		publisher: param.Publisher,
	}
}

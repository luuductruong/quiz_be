package grpchandler

import (
	"context"
	"github.com/quiz_be/services/core/application/job/dto"
	"github.com/quiz_be/services/core/application/job/service"
	appContext "github.com/quiz_be/services/core/context"
	"github.com/quiz_be/services/core/domain/job"
)

type handler struct {
	service.UnimplementedJobServiceServer
	jobDomain job.Service
}

func NewHandler(jobDomain job.Service) service.JobServiceServer {
	return &handler{
		jobDomain: jobDomain,
	}
}

func (h *handler) PushJob(ctx context.Context, req *dto.PushJobReq) (*dto.PushJobResp, error) {
	appCtx := appContext.FromContext(ctx)
	err := h.jobDomain.PushJob(appCtx, req.Topic, req.Data)
	if err != nil {
		return nil, err
	}
	return &dto.PushJobResp{}, nil
}

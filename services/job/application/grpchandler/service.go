package grpchandler

import (
	"github.com/quiz_be/services/core/application/job/service"
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

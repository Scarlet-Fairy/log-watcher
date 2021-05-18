package service

import (
	"context"
	"github.com/go-kit/kit/log"
)

type Service interface {
	QueryImageBuild(ctx context.Context, jobId string, offset, size uint32) ([]Log, error)
	QueryWorkload(ctx context.Context, jobId string, offset, size uint32) ([]Log, error)
}

type basicService struct {
	repository Repository
}

func New(repository Repository, logger log.Logger) Service {
	var service Service
	{
		service = &basicService{
			repository: repository,
		}
		service = LoggingMiddleware(logger)(service)
	}

	return service
}

func (s basicService) QueryImageBuild(ctx context.Context, jobId string, offset, size uint32) ([]Log, error) {
	return s.repository.QueryImageBuild(ctx, jobId, offset, size)
}

func (s basicService) QueryWorkload(ctx context.Context, jobId string, offset, size uint32) ([]Log, error) {
	return s.repository.QueryWorkload(ctx, jobId, offset, size)
}

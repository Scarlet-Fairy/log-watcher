package repository

import (
	"context"
	"github.com/Scarlet-Fairy/log-watcher/pkg/service"
	"github.com/go-kit/kit/log"
)

type Middleware func(repository service.Repository) service.Repository

type loggingMiddleware struct {
	next   service.Repository
	logger log.Logger
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(repository service.Repository) service.Repository {
		return &loggingMiddleware{
			next:   repository,
			logger: logger,
		}
	}
}

func (l loggingMiddleware) QueryImageBuild(ctx context.Context, jobId string, offset, size uint32) (logs []service.Log, err error) {
	defer func() {
		l.logger.Log(
			"method", "QueryImageBuild",
			"jobId", jobId,
			"offset", offset,
			"size", size,
			"logs", logs,
			"err", err,
		)
	}()

	return l.next.QueryImageBuild(ctx, jobId, offset, size)
}

func (l loggingMiddleware) QueryWorkload(ctx context.Context, jobId string, offset, size uint32) (logs []service.Log, err error) {
	defer func() {
		l.logger.Log(
			"method", "QueryWorkload",
			"jobId", jobId,
			"offset", offset,
			"size", size,
			"logs", logs,
			"err", err,
		)
	}()

	return l.next.QueryWorkload(ctx, jobId, offset, size)
}

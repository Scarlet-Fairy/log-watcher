package service

import (
	"context"
	"github.com/go-kit/kit/log"
)

type Middleware func(service Service) Service

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(service Service) Service {
		return &loggingMiddleware{
			next:   service,
			logger: logger,
		}
	}
}

func (l loggingMiddleware) QueryImageBuild(ctx context.Context, jobId string, offset, size uint32) (logs []Log, err error) {
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

func (l loggingMiddleware) QueryWorkload(ctx context.Context, jobId string, offset, size uint32) (logs []Log, err error) {
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

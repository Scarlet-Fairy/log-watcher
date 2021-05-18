package endpoint

import (
	"context"
	"github.com/Scarlet-Fairy/log-watcher/pkg/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type LogWatcherEndoint struct {
	GetLogsEndpoint endpoint.Endpoint
}

func NewEndpoint(s service.Service, logger log.Logger) LogWatcherEndoint {
	var getLogsEndpoint endpoint.Endpoint
	{
		getLogsEndpoint = makeGetLogsEndpoint(s)
		getLogsEndpoint = LoggingMiddleware(log.With(logger, "method", "GetLogs"))(getLogsEndpoint)
		getLogsEndpoint = UnwrapErrorMiddleware()(getLogsEndpoint)
	}

	return LogWatcherEndoint{
		GetLogsEndpoint: getLogsEndpoint,
	}
}

var (
	_ endpoint.Failer = GetLogsResponse{}
)

type GetLogsRequest struct {
	JobId  string
	Offset uint32
	Size   uint32
}

type GetLogsResponse struct {
	Logs []service.Log
	Err  error `json:"-"`
}

func (r GetLogsResponse) Failed() error {
	return r.Err
}

func makeGetLogsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*GetLogsRequest)

		logs, err := s.QueryWorkload(
			ctx,
			req.JobId,
			req.Offset,
			req.Size,
		)

		return &GetLogsResponse{
			Logs: logs,
			Err:  err,
		}, nil
	}
}

package grpc

import (
	"context"
	"github.com/Scarlet-Fairy/log-watcher/pb"
	"github.com/Scarlet-Fairy/log-watcher/pkg/endpoint"
)

func decodeGetLogsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetLogsRequest)

	return &endpoint.GetLogsRequest{
		JobId:  req.DeployId,
		Offset: req.Offset,
		Size:   req.Size,
	}, nil
}

func encodeDeployResponse(_ context.Context, resp interface{}) (interface{}, error) {
	res := resp.(*endpoint.GetLogsResponse)

	var logs []*pb.Log
	for _, log := range res.Logs {
		logs = append(logs, &pb.Log{
			Timestamp: log.Timestamp.String(),
			Body:      log.Message,
		})
	}
	return &pb.GetLogsResponse{
		Logs: logs,
	}, nil
}

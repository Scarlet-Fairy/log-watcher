package grpc

import (
	"context"
	"github.com/Scarlet-Fairy/log-watcher/pb"
	"github.com/Scarlet-Fairy/log-watcher/pkg/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	pb.UnimplementedLogWatcherServer
	getLogs grpctransport.Handler
}

func NewGRPCServer(endpoints endpoint.LogWatcherEndoint, logger log.Logger) pb.LogWatcherServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		getLogs: grpctransport.NewServer(
			endpoints.GetLogsEndpoint,
			decodeGetLogsRequest,
			encodeDeployResponse,
			options...,
		),
	}
}

func (g grpcServer) GetLogs(ctx context.Context, request *pb.GetLogsRequest) (*pb.GetLogsResponse, error) {
	_, resp, err := g.getLogs.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.GetLogsResponse), nil
}

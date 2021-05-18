package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Scarlet-Fairy/log-watcher/pb"
	"github.com/Scarlet-Fairy/log-watcher/pkg/endpoint"
	"github.com/Scarlet-Fairy/log-watcher/pkg/logger"
	elasticSearchRepository "github.com/Scarlet-Fairy/log-watcher/pkg/repository/elasticsearch"
	"github.com/Scarlet-Fairy/log-watcher/pkg/service"
	grpcTransport "github.com/Scarlet-Fairy/log-watcher/pkg/transport/grpc"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-kit/kit/log/level"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/run"
	"google.golang.org/grpc"
	"net"
	"os"
)

var (
	grpcAddr   = flag.String("grpc-url", ":8084", "gRPC server listen address")
	elasticUrl = flag.String("es-url", "http://localhost:9200", "url of elasticsearch instance")
)

var (
	loggers                   = logger.NewLogger()
	infoLogger                = loggers.InfoLogger
	warnLogger                = loggers.WarnLogger
	errorLogger               = loggers.ErrorLogger
	transportLayerLogger      = loggers.TransportLayerLogger
	endpointLayerLogger       = loggers.EndpointLayerLogger
	serviceComponentLogger    = loggers.ServiceComponentLogger
	repositoryComponentLogger = loggers.RepositoryComponentLogger
)

var ctx = context.Background()

func main() {
	flag.Parse()

	elasticSearchClient, err := newElasticSearchClient(*elasticUrl)
	if err != nil {
		errorLogger.Log(
			"elastic-url", *elasticUrl,
			"msg", "failed to init elastic client",
			"err", err,
		)
	}

	repositoryInstance := elasticSearchRepository.New(elasticSearchClient, repositoryComponentLogger)

	svc := service.New(repositoryInstance, serviceComponentLogger)
	endpoints := endpoint.NewEndpoint(svc, endpointLayerLogger)
	grpcServer := grpcTransport.NewGRPCServer(endpoints, transportLayerLogger)

	var g run.Group
	{
		grpcListener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			level.Error(transportLayerLogger).Log(
				"during", "init",
				"msg", fmt.Sprintf("failed to listen on %s", *grpcAddr),
				"err", err,
			)
			os.Exit(1)
		}

		g.Add(func() error {
			transportLayerLogger.Log(
				"addr", *grpcAddr,
			)

			baseServer := grpc.NewServer(
				grpc.UnaryInterceptor(kitgrpc.Interceptor),
			)
			pb.RegisterLogWatcherServer(baseServer, grpcServer)

			return baseServer.Serve(grpcListener)
		}, func(err error) {
			if err = grpcListener.Close(); err != nil {
				panic(err)
			}
		})
	}

	infoLogger.Log("exit", g.Run())
}

func newElasticSearchClient(url string) (*elasticsearch.Client, error) {
	return elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			url,
		},
	})
}

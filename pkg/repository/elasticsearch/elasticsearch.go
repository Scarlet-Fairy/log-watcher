package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Scarlet-Fairy/log-watcher/pkg/repository"
	"github.com/Scarlet-Fairy/log-watcher/pkg/service"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"time"
)

type elasticsearchRepository struct {
	client *elasticsearch.Client
}

const (
	CoreIndex       = "scarlet-fairy-core"
	WorkloadIndex   = "scarlet-fairy-workload"
	ImageBuildIndex = "scarlet-fairy-imagebuild"
)

func New(client *elasticsearch.Client, logger log.Logger) service.Repository {
	var instance service.Repository
	{
		instance = &elasticsearchRepository{
			client: client,
		}
		instance = repository.LoggingMiddleware(logger)(instance)
	}

	return instance
}

func (e elasticsearchRepository) query(
	ctx context.Context,
	index string,
	jobId string,
	from uint32,
	size uint32,
) ([]service.Log, error) {
	queryBody := ReadLogsQuery{
		Query: Query{
			Match: map[string]interface{}{
				"agent.name": jobId,
			},
		},
	}
	marshaledBody, err := json.Marshal(queryBody)
	if err != nil {
		return nil, err
	}

	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex(index),
		e.client.Search.WithBody(bytes.NewReader(marshaledBody)),
		e.client.Search.WithFrom(int(from)),
		e.client.Search.WithSize(int(size)),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	if res.IsError() {
		var resError map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&resError); err != nil {
			return nil, err
		}

		return nil, errors.New(
			fmt.Sprintf(
				"[%s] Type: %s, Reason: %s",
				res.Status(),
				resError["type"],
				resError["reason"],
			),
		)
	}

	var resBody map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&res); err != nil {
		return nil, err
	}

	hits := resBody["hits"].(map[string]interface{})["hits"].([]map[string]interface{})

	var logs []service.Log
	for _, hit := range hits {
		source := hit["_source"].(*SearchResponse)
		timestamp, err := time.Parse(time.RFC3339, source.Timestamp)
		if err != nil {
			return nil, err
		}

		logs = append(logs, service.Log{
			Id:        hit["_id"].(string),
			Timestamp: timestamp,
			Message:   source.Message,
		})
	}

	return logs, nil
}

func (e elasticsearchRepository) QueryImageBuild(ctx context.Context, jobId string, offset, size uint32) ([]service.Log, error) {
	return e.query(
		ctx,
		ImageBuildIndex,
		jobId,
		offset,
		size,
	)
}

func (e elasticsearchRepository) QueryWorkload(ctx context.Context, jobId string, offset, size uint32) ([]service.Log, error) {
	return e.query(
		ctx,
		WorkloadIndex,
		jobId,
		offset,
		size,
	)
}

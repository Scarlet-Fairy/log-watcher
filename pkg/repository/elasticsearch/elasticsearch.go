package elasticsearch

import (
	"context"
	"encoding/json"
	"github.com/Scarlet-Fairy/log-watcher/pkg/repository"
	"github.com/Scarlet-Fairy/log-watcher/pkg/service"
	"github.com/go-kit/kit/log"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"time"
)

type elasticsearchRepository struct {
	client *elastic.Client
}

const (
	CoreIndex       = "scarlet-fairy-core"
	WorkloadIndex   = "scarlet-fairy-workloads"
	ImageBuildIndex = "scarlet-fairy-imagebuild"
)

func New(client *elastic.Client, logger log.Logger) service.Repository {
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

	res, err := e.client.Search().
		Index(index).
		From(int(from)).
		Size(int(size)).
		Sort("@timestamp", false).
		Query(elastic.NewMatchQuery("agent.name", jobId)).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	var logs []service.Log
	for _, hit := range res.Hits.Hits {
		var source map[string]interface{}
		if err := json.Unmarshal(hit.Source, &source); err != nil {
			return nil, errors.Wrap(err, "Unmarshaling source")
		}

		timestamp, err := time.Parse(time.RFC3339, source["@timestamp"].(string))
		if err != nil {
			return nil, errors.Wrap(err, "Parsing timestamp")
		}

		logs = append(logs, service.Log{
			Id:        hit.Id,
			Timestamp: timestamp,
			Message:   source["message"].(string),
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

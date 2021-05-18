package service

import "context"

type Repository interface {
	QueryImageBuild(ctx context.Context, jobId string, offset, size uint32) ([]Log, error)
	QueryWorkload(ctx context.Context, jobId string, offset, size uint32) ([]Log, error)
}

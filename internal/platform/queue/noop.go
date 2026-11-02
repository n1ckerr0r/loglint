package queue

import (
	"context"

	"github.com/n1ckerr0r/loglint/internal/service"
)

type NoopQueue struct{}

func (NoopQueue) Publish(context.Context, service.Job) error {
	return nil
}

func (NoopQueue) Consume(context.Context) (<-chan service.Job, error) {
	jobs := make(chan service.Job)
	close(jobs)
	return jobs, nil
}

func (NoopQueue) Close() error {
	return nil
}

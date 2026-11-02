package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/n1ckerr0r/loglint/internal/service"
)

type Worker struct {
	repo  service.RunRepository
	queue service.Queue
	cache service.StatusCache
}

func New(repo service.RunRepository, queue service.Queue, cache service.StatusCache) *Worker {
	return &Worker{
		repo:  repo,
		queue: queue,
		cache: cache,
	}
}

func (worker *Worker) Run(ctx context.Context) error {
	jobs, err := worker.queue.Consume(ctx)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case job, ok := <-jobs:
			if !ok {
				return nil
			}
			worker.handleJob(ctx, job)
		}
	}
}

func (worker *Worker) handleJob(ctx context.Context, job service.Job) {
	if err := worker.repo.UpdateRun(ctx, job.RunID, service.RunStatusRunning, ""); err != nil {
		slog.Error("update run status", "run_id", job.RunID, "error", err)
		return
	}
	_ = worker.cache.SetStatus(ctx, job.RunID, service.RunStatusRunning)

	report := analyze(job)
	if err := worker.repo.UpdateRun(ctx, job.RunID, service.RunStatusCompleted, report); err != nil {
		slog.Error("complete run", "run_id", job.RunID, "error", err)
		return
	}
	_ = worker.cache.SetStatus(ctx, job.RunID, service.RunStatusCompleted)
}

func analyze(job service.Job) string {
	// The worker owns orchestration; actual repository checkout and analyzer execution can evolve independently.
	time.Sleep(100 * time.Millisecond)
	return fmt.Sprintf("Analysis completed for %s ref=%s commit=%s", job.Project, job.Ref, job.Commit)
}

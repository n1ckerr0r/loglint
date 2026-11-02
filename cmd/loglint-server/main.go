package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	platformconfig "github.com/n1ckerr0r/loglint/internal/platform/config"
	"github.com/n1ckerr0r/loglint/internal/platform/queue"
	"github.com/n1ckerr0r/loglint/internal/platform/storage"
	"github.com/n1ckerr0r/loglint/internal/server"
	"github.com/n1ckerr0r/loglint/internal/service"
)

func main() {
	cfg := platformconfig.FromEnv()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	repo, err := repository(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("repository: %v", err)
	}
	defer closeResource("repository", repo)

	jobQueue, err := jobQueue(cfg.RabbitMQURL, cfg.QueueName)
	if err != nil {
		log.Fatalf("queue: %v", err)
	}
	defer closeResource("queue", jobQueue)

	app := server.New(repo, jobQueue)
	log.Printf("loglint api listening on %s", cfg.Addr)
	if err := server.ListenAndServe(ctx, cfg.Addr, app.Handler(), cfg.Shutdown); err != nil {
		log.Fatal(err)
	}
}

func repository(ctx context.Context, databaseURL string) (service.RunRepository, error) {
	if databaseURL == "" {
		return storage.NewMemoryRepository(), nil
	}
	return storage.NewPostgresRepository(ctx, databaseURL)
}

func jobQueue(rabbitURL, queueName string) (service.Queue, error) {
	if rabbitURL == "" {
		return queue.NoopQueue{}, nil
	}
	return queue.NewRabbitMQQueue(rabbitURL, queueName)
}

func closeResource(name string, closer interface{ Close() error }) {
	if err := closer.Close(); err != nil {
		log.Printf("close %s: %v", name, err)
	}
}

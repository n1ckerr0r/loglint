package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/n1ckerr0r/loglint/internal/platform/cache"
	platformconfig "github.com/n1ckerr0r/loglint/internal/platform/config"
	"github.com/n1ckerr0r/loglint/internal/platform/queue"
	"github.com/n1ckerr0r/loglint/internal/platform/storage"
	"github.com/n1ckerr0r/loglint/internal/service"
	"github.com/n1ckerr0r/loglint/internal/worker"
)

func main() {
	cfg := platformconfig.FromEnv()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if cfg.DatabaseURL == "" || cfg.RabbitMQURL == "" {
		log.Fatal("DATABASE_URL and RABBITMQ_URL are required for worker")
	}

	repo, err := storage.NewPostgresRepository(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("repository: %v", err)
	}
	defer closeResource("repository", repo)

	jobQueue, err := queue.NewRabbitMQQueue(cfg.RabbitMQURL, cfg.QueueName)
	if err != nil {
		log.Fatalf("queue: %v", err)
	}
	defer closeResource("queue", jobQueue)

	statusCache, err := statusCache(cfg.RedisURL)
	if err != nil {
		log.Fatalf("cache: %v", err)
	}
	defer closeResource("cache", statusCache)

	log.Println("loglint worker started")
	if err := worker.New(repo, jobQueue, statusCache).Run(ctx); err != nil && ctx.Err() == nil {
		log.Fatal(err)
	}
}

func statusCache(redisURL string) (service.StatusCache, error) {
	if redisURL == "" {
		return cache.NoopStatusCache{}, nil
	}
	return cache.NewRedisStatusCache(redisURL)
}

func closeResource(name string, closer interface{ Close() error }) {
	if err := closer.Close(); err != nil {
		log.Printf("close %s: %v", name, err)
	}
}

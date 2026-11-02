package queue

import (
	"context"
	"encoding/json"

	"github.com/n1ckerr0r/loglint/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQQueue struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

func NewRabbitMQQueue(rabbitURL, queueName string) (*RabbitMQQueue, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	if _, err := channel.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return nil, err
	}
	return &RabbitMQQueue{
		conn:      conn,
		channel:   channel,
		queueName: queueName,
	}, nil
}

func (queue *RabbitMQQueue) Publish(ctx context.Context, job service.Job) error {
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return queue.channel.PublishWithContext(
		ctx,
		"",
		queue.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

func (queue *RabbitMQQueue) Consume(ctx context.Context) (<-chan service.Job, error) {
	deliveries, err := queue.channel.Consume(queue.queueName, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	jobs := make(chan service.Job)
	go func() {
		defer close(jobs)
		for {
			select {
			case <-ctx.Done():
				return
			case delivery, ok := <-deliveries:
				if !ok {
					return
				}
				var job service.Job
				if err := json.Unmarshal(delivery.Body, &job); err != nil {
					_ = delivery.Nack(false, false)
					continue
				}
				jobs <- job
				_ = delivery.Ack(false)
			}
		}
	}()
	return jobs, nil
}

func (queue *RabbitMQQueue) Close() error {
	if err := queue.channel.Close(); err != nil {
		_ = queue.conn.Close()
		return err
	}
	return queue.conn.Close()
}

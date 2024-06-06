package saga

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func nackWithDelay(msg *amqp.Delivery, queueName string, delay time.Duration, maxRetries int) (int, error) {
	channel, err := getConsumeChannel()
	if err != nil {
		return 0, fmt.Errorf("error getting consume channel in nackWithDelay: %w", err)
	}
	err = channel.Nack(msg.DeliveryTag, false, false) // nack without requeueing immediately
	if err != nil {
		return 0, fmt.Errorf("error nacking message: %w", err)
	}

	count := 1
	if countHeader, ok := msg.Headers["x-death"]; ok {
		if deaths, ok := countHeader.([]interface{}); ok && len(deaths) > 0 {
			if death, ok := deaths[0].(amqp.Table); ok {
				if deathCount, ok := death["count"].(int64); ok {
					count = int(deathCount) + 1
				}
			}
		}
	}

	if count > maxRetries {
		fmt.Printf("MAX NACK RETRIES REACHED: %d - NACKING %s - %s", maxRetries, queueName, msg.Body)
		return maxRetries, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = channel.PublishWithContext(
		ctx,
		string(RequeueExchange),
		fmt.Sprintf("%s_routing_key", queueName),
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Expiration: fmt.Sprintf("%d", delay.Milliseconds()),
			Headers:    msg.Headers,
			Body:       msg.Body,
		},
	)
	if err != nil {
		return 0, fmt.Errorf("error publishing message: %w", err)
	}
	return count, nil
}

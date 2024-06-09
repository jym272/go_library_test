package saga

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumeChannel struct {
	channel   *amqp.Channel
	msg       *amqp.Delivery
	queueName string
}

const (
	NACKING_DELAY_MS = 5000 // 5 seconds
	MAX_NACK_RETRIES = 3
	// MAX_OCCURRENCE
	/**
	 * Define the maximum occurrence in a fail saga step of the nack delay with fibonacci strategy
	 * | Occurrence | Delay in the next nack |
	 * |------------|------------------------|
	 * | 17         | 0.44 hours  |
	 * | 18         | 0.72 hours  |
	 * | 19         | 1.18 hours  |
	 * | 20         | 1.88 hours  |
	 * | 21         | 3.04 hours  |
	 * | 22         | 4.92 hours  |
	 * | 23         | 7.96 hours  |
	 * | 24         | 12.87 hours |
	 * | 25         | 20.84 hours |
	 */
	MAX_OCCURRENCE = 19
)

// math.MaxInt8
func (c *ConsumeChannel) NackWithDelay(delay time.Duration, maxRetries int) (int, time.Duration, error) {
	err := c.channel.Nack(c.msg.DeliveryTag, false, false)
	if err != nil {
		return 0, 0, fmt.Errorf("error nacking message: %w", err)
	}

	count := 0
	if retryCount, ok := c.msg.Headers["x-retry-count"]; ok {
		count = retryCount.(int)
	}
	count++

	if count > maxRetries {
		fmt.Printf("MAX NACK RETRIES REACHED: %d - NACKING %s - %s", maxRetries, c.queueName, c.msg.Body)
		return count, delay, nil
	}

	c.msg.Headers["x-retry-count"] = count
	err = c.publishNackEvent(delay)
	if err != nil {
		return 0, 0, fmt.Errorf("error publishing nack event: %w", err)
	}
	return count, delay, nil
}
func (c *ConsumeChannel) NackWithFibonacciStrategy(maxOccurrence, maxRetries int) (int, time.Duration, int, error) {
	err := c.channel.Nack(c.msg.DeliveryTag, false, false)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error nacking message: %w", err)
	}

	count := 0
	if retryCount, ok := c.msg.Headers["x-retry-count"]; ok {
		count = retryCount.(int)
	}
	count++

	occurrence := 0
	if o, ok := c.msg.Headers["x-occurrence"]; ok {
		if o.(int) >= maxOccurrence {
			// the occurrence is reset to 0 to avoid large delay in the next nack
			occurrence = 0
		}
	}
	occurrence++

	delay := time.Duration(fibonacci(occurrence)) * time.Second

	if count > maxRetries {
		fmt.Printf("MAX NACK RETRIES REACHED: %d - NACKING %s - %s", maxRetries, c.queueName, c.msg.Body)
		return count, delay, occurrence, nil
	}

	c.msg.Headers["x-retry-count"] = count
	c.msg.Headers["x-occurrence"] = occurrence

	err = c.publishNackEvent(delay)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error publishing nack event: %w", err)
	}
	return count, delay, occurrence, nil
}

func (c *ConsumeChannel) publishNackEvent(delay time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if c.msg.Exchange == string(MatchingExchange) {
		// the header that is deleted is the one that has all the micros listening to a certain event,
		// otherwise the nacking reaches everyone.
		delete(c.msg.Headers, "all-micro")
		// deliver to one micro in particular, the one that is nacking
		c.msg.Headers["micro"] = c.queueName
		err := c.channel.PublishWithContext(
			ctx,
			string(MatchingRequeueExchange),
			"",
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				Expiration:   fmt.Sprintf("%d", delay.Milliseconds()),
				Headers:      c.msg.Headers,
				Body:         c.msg.Body,
				DeliveryMode: amqp.Persistent,
			},
		)
		if err != nil {
			return fmt.Errorf("error publishing message: %w", err)
		}

	} else {
		// is a saga event
		err := c.channel.PublishWithContext(
			ctx,
			string(RequeueExchange),
			fmt.Sprintf("%s_routing_key", c.queueName),
			false,
			false,
			amqp.Publishing{
				Expiration:   fmt.Sprintf("%d", delay.Milliseconds()),
				Headers:      c.msg.Headers,
				Body:         c.msg.Body,
				DeliveryMode: amqp.Persistent,
			},
		)
		if err != nil {
			return fmt.Errorf("error publishing message: %w", err)
		}
	}
	return nil
}

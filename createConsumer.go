package saga

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// https://blog.rabbitmq.com/posts/2012/04/rabbitmq-performance-measurements-part-2/
func createConsumers(consumers []QueueConsumerProps) error {
	channel, err := getConsumeChannel()
	if err != nil {
		return err
	}

	for _, consumer := range consumers {
		queueName := consumer.QueueName
		exchange := string(consumer.Exchange)
		requeueQueue := fmt.Sprintf("%s_requeue", queueName)
		routingKey := fmt.Sprintf("%s_routing_key", queueName)

		// Assert exchange and queue for the consumer.
		err := channel.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
		if err != nil {
			return err
		}
		_, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
		if err != nil {
			return err
		}
		err = channel.QueueBind(queueName, routingKey, exchange, false, nil)
		if err != nil {
			return err
		}

		// Set up requeue mechanism by creating a requeue exchange and binding requeue queue to it.
		err = channel.ExchangeDeclare(string(RequeueExchange), "direct", true, false, false, false, nil)
		if err != nil {
			return err
		}
		_, err = channel.QueueDeclare(requeueQueue, true, false, false, false, amqp.Table{
			"x-dead-letter-exchange": exchange,
		})
		if err != nil {
			return err
		}
		err = channel.QueueBind(requeueQueue, routingKey, string(RequeueExchange), false, nil)
		if err != nil {
			return err
		}

		// Set the prefetch count to process only one message at a time to maintain order and control concurrency.
		// TODO: is the same channel as the createHeaderConsumers channel, if for some reason the prefetch count is changed in one place, it will be changed in the other
		// TODO: solution, create a new channel
		err = channel.Qos(1, 0, false)
		if err != nil {
			return err
		}
	}

	return nil
}

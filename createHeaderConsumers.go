package saga

import (
	"fmt"
	"slices"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

func getEventKey(event MicroserviceEvent) string {
	return strings.ToUpper(string(event))
}

func getEventObject(event MicroserviceEvent) amqp.Table {
	key := getEventKey(event)
	return amqp.Table{key: event}
}

func createHeaderConsumers(queueName string, events []MicroserviceEvent) error {
	channel, err := getConsumeChannel()
	if err != nil {
		return fmt.Errorf("failed to get consume channel: %w", err)
	}

	requeueQueue := fmt.Sprintf("%s_matching_requeue", queueName)

	for _, exchange := range []Exchange{MatchingExchange, MatchingRequeueExchange} {
		err = channel.ExchangeDeclare(string(exchange), "headers", true, false, false, false, nil)
		if err != nil {
			return fmt.Errorf("failed to declare exchange %s: %w", string(exchange), err)
		}
	}

	_, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}
	_, err = channel.QueueDeclare(requeueQueue, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange": string(MatchingExchange),
	})
	if err != nil {
		return fmt.Errorf("failed to declare requeue queue %s: %w", requeueQueue, err)
	}

	// Handle individual events
	for _, ev := range microserviceEventValues() {
		headerEvent := getEventObject(ev)

		// Declare and bind event-specific exchanges
		err = channel.ExchangeDeclare(string(ev), "headers", true, false, false, false, nil)
		if err != nil {
			return fmt.Errorf("failed to declare exchange %s: %w", ev, err)
		}
		headersArgs := amqp.Table{
			"all-micro": "yes",
			"x-match":   "all",
		}
		for k, v := range headerEvent {
			headersArgs[k] = v
		}
		err = channel.ExchangeBind(string(ev), "", string(MatchingExchange), false, headersArgs)
		if err != nil {
			return fmt.Errorf("failed to bind exchange %s to %s: %w", ev, MatchingExchange, err)
		}

		requeueExchange := fmt.Sprintf("%s_requeue", ev)

		err = channel.ExchangeDeclare(requeueExchange, "headers", true, false, false, false, nil)
		if err != nil {
			return fmt.Errorf("failed to declare requeue exchange %s: %w", requeueExchange, err)
		}
		err = channel.ExchangeBind(requeueExchange, "", string(MatchingRequeueExchange), false, headerEvent)
		if err != nil {
			return fmt.Errorf("failed to bind requeue exchange %s to %s: %w", requeueExchange, MatchingRequeueExchange, err)
		}

		headersArgs = amqp.Table{
			"micro":   queueName,
			"x-match": "all",
		}
		for k, v := range headerEvent {
			headersArgs[k] = v
		}
		if slices.Contains(events, ev) {
			// Bindings for included events
			err = channel.QueueBind(queueName, "", string(ev), false, headerEvent)
			if err != nil {
				return fmt.Errorf("failed to bind queue %s to exchange %s: %w", queueName, ev, err)
			}

			err = channel.QueueBind(requeueQueue, "", fmt.Sprintf("%s_requeue", ev), false, headersArgs)
			if err != nil {
				return fmt.Errorf("failed to bind requeue queue %s to exchange %s_requeue: %w", requeueQueue, ev, err)
			}

			microExchange := fmt.Sprintf("%s_%s", ev, queueName)
			err = channel.ExchangeDeclare(microExchange, "headers", true, false, false, false, nil)
			if err != nil {
				return fmt.Errorf("failed to declare exchange %s: %w", microExchange, err)
			}
			err = channel.ExchangeBind(microExchange, "", string(MatchingExchange), false, headersArgs)
			if err != nil {
				return fmt.Errorf("failed to bind exchange %s to %s: %w", microExchange, MatchingExchange, err)
			}

			err = channel.QueueBind(queueName, "", microExchange, false, headersArgs)
			if err != nil {
				return fmt.Errorf("failed to bind queue %s to exchange %s: %w", queueName, microExchange, err)
			}

		} else {
			// Attempt to unbind the queue, ignoring errors if it's already unbound
			err = channel.QueueUnbind(queueName, "", string(ev), headerEvent)
			if err != nil {
				return fmt.Errorf("failed to unbind queue %s from exchange %s: %w", queueName, ev, err)
			}
			err = channel.QueueUnbind(requeueQueue, "", fmt.Sprintf("%s_requeue", ev), headersArgs)
			if err != nil {
				return fmt.Errorf("failed to unbind requeue queue %s from exchange %s_requeue: %w", requeueQueue, ev, err)
			}
			err = channel.ExchangeDelete(fmt.Sprintf("%s_%s", ev, queueName), false, false)
			if err != nil {
				return fmt.Errorf("failed to delete exchange %s_%s: %w", ev, queueName, err)
			}
		}
	}
	// TODO: is the same channel as the createConsumers channel, if for some reason the prefetch count is changed in one place, it will be changed in the other
	// TODO: solution, create a new channel
	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}
	return nil
}

func containsEvent(events []MicroserviceEvent, ev MicroserviceEvent) bool {
	for _, e := range events {
		if e == ev {
			return true
		}
	}
	return false

}

package saga

import (
	"fmt"
	"github.com/jym272/go_library_test/event"
	"github.com/jym272/go_library_test/micro"
	amqp "github.com/rabbitmq/amqp091-go"
)

// consume consumes messages from the queue and processes them.
func consume[T any, U comparable](e *Emitter[T, U], queueName string, cb func(*amqp.Delivery, *amqp.Channel, *Emitter[T, U], string)) error {
	channel, err := getConsumeChannel()
	if err != nil {
		return fmt.Errorf("failed to get channel: %w", err)
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

	channelQ, err := channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	for msg := range channelQ {
		cb(&msg, channel, e, queueName)
	}

	return nil
}

func getQueueName(microservice micro.AvailableMicroservices) string {
	return fmt.Sprintf("%s_saga_commands", microservice)
}

func getQueueConsumer(microservice micro.AvailableMicroservices) QueueConsumerProps {
	return QueueConsumerProps{
		QueueName: getQueueName(microservice),
		Exchange:  CommandsExchange,
	}
}

type Transactional struct {
	RabbitUri    string
	Microservice micro.AvailableMicroservices
	Events       []event.MicroserviceEvent
	isReady      bool
}

func (t *Transactional) ConnectToSagaCommandEmitter() *Emitter[CommandHandler, micro.StepCommand] {
	t.prepare()
	q := getQueueConsumer(t.Microservice)
	e := newEmitter[CommandHandler, micro.StepCommand]()

	err := createConsumers([]QueueConsumerProps{q})
	if err != nil {
		panic(err)
	}

	go func() {
		err = consume(e, q.QueueName, sagaCommandCallback)
		if err != nil {
			fmt.Println("Error consuming messages:", err)
		}
	}()

	return e
}

func (t *Transactional) ConnectToEvents() *Emitter[EventHandler, event.MicroserviceEvent] {
	t.prepare()
	queueName := fmt.Sprintf("%s_match_commands", t.Microservice)
	e := newEmitter[EventHandler, event.MicroserviceEvent]()

	err := createHeaderConsumer(queueName, t.Events)
	if err != nil {
		panic(err)
	}

	go func() {
		err = consume(e, queueName, eventCallback)
		if err != nil {
			fmt.Println("Error consuming messages:", err)
		}
	}()

	return e
}

/// some commasdasdentasdasdasd

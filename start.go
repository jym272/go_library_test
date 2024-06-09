package saga

import (
	"fmt"
)

type commandEmitterConf struct {
	// microservice is the microservice that will be connecting to the events.
	microservice AvailableMicroservices
}

func getQueueName(microservice AvailableMicroservices) string {
	return fmt.Sprintf("%s_saga_commands", microservice)
}

func getQueueConsumer(microservice AvailableMicroservices) QueueConsumerProps {
	return QueueConsumerProps{
		QueueName: getQueueName(microservice),
		Exchange:  CommandsExchange,
	}
}

func connectToSagaCommandEmitter(conf commandEmitterConf) (*Emitter[CommandHandler, StepCommand], error) {
	q := getQueueConsumer(conf.microservice)
	e := newEmitter[CommandHandler, StepCommand]()

	err := createConsumers([]QueueConsumerProps{q})
	if err != nil {
		return nil, err
	}

	go func() {
		err := consume(e, q.QueueName, microserviceConsumeCallback)
		if err != nil {
			fmt.Println("Error consuming messages:", err)
		}
	}()

	return e, nil
}

type eventsConf struct {
	// microservice is the microservice that will be connecting to the events.
	microservice AvailableMicroservices
	// events is the list of events that the microservice will be connecting to.
	events []MicroserviceEvent
}

func connectToEvents(conf eventsConf) (*Emitter[EventHandler, MicroserviceEvent], error) {

	queueName := fmt.Sprintf("%s_match_commands", conf.microservice)
	e := newEmitter[EventHandler, MicroserviceEvent]()

	err := createHeaderConsumer(queueName, conf.events)
	if err != nil {
		return nil, fmt.Errorf("error creating header consumers: %w", err)
	}

	go func() {
		err = consume(e, queueName, eventCallback)
		if err != nil {
			fmt.Println("Error consuming messages:", err)
		}
	}()

	return e, nil
}

type TransactionalConfig struct {
	Url          string
	Microservice AvailableMicroservices
	Events       []MicroserviceEvent
}

func StartTransactional(config TransactionalConfig) (*Emitter[EventHandler, MicroserviceEvent], *Emitter[CommandHandler, StepCommand], error) {
	if err := Prepare(config.Url); err != nil {
		return nil, nil, fmt.Errorf("error preparing transactional: %w", err)
	}

	eventEmitter, err := connectToEvents(eventsConf{
		microservice: config.Microservice,
		events:       config.Events,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to events: %w", err)
	}

	commandEmitter, err := connectToSagaCommandEmitter(commandEmitterConf{
		microservice: config.Microservice,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to command emitter: %w", err)

	}

	return eventEmitter, commandEmitter, nil
}

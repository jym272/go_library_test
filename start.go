package saga

import (
	"fmt"
)

type CommandEmitterConf struct {
	// microservice is the microservice that will be connecting to the events.
	microservice AvailableMicroservices
}

func connectToSagaCommandEmitter(conf CommandEmitterConf) (*Emitter[CommandHandler, SagaStepCommand], error) {
	q := getQueueConsumer(conf.microservice)
	e := NewEmitter[CommandHandler, SagaStepCommand]()

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

type EventsConf struct {
	// microservice is the microservice that will be connecting to the events.
	microservice AvailableMicroservices
	// events is the list of events that the microservice will be connecting to.
	events []MicroserviceEvent
}

func connectToEvents(conf EventsConf) (*Emitter[EventHandler, MicroserviceEvent], error) {

	q := fmt.Sprintf("%s_match_commands", conf.microservice)
	e := NewEmitter[EventHandler, MicroserviceEvent]()

	err := createHeaderConsumers(q, conf.events)
	if err != nil {
		return nil, err
	}

	go func() {
		err := consume(e, q, eventCallback)
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

func StartTransactional(config TransactionalConfig) (*Emitter[EventHandler, MicroserviceEvent], *Emitter[CommandHandler, SagaStepCommand], error) {
	if err := prepare(config.Url); err != nil {
		return nil, nil, fmt.Errorf("error preparing transactional: %w", err)
	}

	eventEmitter, err := connectToEvents(EventsConf{
		microservice: config.Microservice,
		events:       config.Events,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to events: %w", err)
	}

	commandEmitter, err := connectToSagaCommandEmitter(CommandEmitterConf{
		microservice: config.Microservice,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to command emitter: %w", err)

	}

	return eventEmitter, commandEmitter, nil
}

func Foo() {

	eventEmitter, commandEmitter, err := StartTransactional(TransactionalConfig{
		Url:          "amqp://guest:guest@localhost:5672/",
		Microservice: Auth,
		Events: []MicroserviceEvent{
			SocialNewUserEvent,
		},
	})
	if err != nil {
		fmt.Println("Error starting transactional:", err)
		return
	}
	eventEmitter.On(SocialNewUserEvent, func(handler EventHandler) {
		payload := handler.Payload.(SocialNewUserPayload)
		fmt.Println("SocialNewUserEvent received", payload)
	})

	commandEmitter.On(NewUserSetRolesToRoomsCommand, func(handler CommandHandler) {
		payload := handler.Payload
		fmt.Println("SocialNewUserEvent received", payload)
	})

}

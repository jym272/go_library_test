package saga

import (
	"fmt"
)

// dummy comment.
func ConnectToSagaCommandEmitter(url string, microservice AvailableMicroservices) (*Emitter, error) {
	err := prepare(url)
	if err != nil {
		return nil, err
	}
	q := getQueueConsumer(microservice)
	e := NewEmitter()

	err = createConsumers([]QueueConsumerProps{q})
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

type EventsConfig struct {
	// url is the URL of the RabbitMQ server.
	url string
	// microservice is the microservice that will be connecting to the events.
	microservice AvailableMicroservices
	// events is the list of events that the microservice will be connecting to.
	events []MicroserviceEvent
}

func connectToEvents(conf EventsConfig) (*Emitter, error) {

	err := prepare(conf.url)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf("%s_match_commands", conf.microservice)
	e := NewEmitter()

	err = createHeaderConsumers(q, conf.events)
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
func foo() {
	e, _ := connectToEvents(EventsConfig{
		url:          "url",
		microservice: "microservice",
		events: []MicroserviceEvent{
			"event1",
			"event2",
		},
	})
	e.On("", func(CommandHandler) {})
}

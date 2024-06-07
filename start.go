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

func connectToEvents(url string, microservice AvailableMicroservices, events []MicroserviceEvent) (*Emitter, error) {
	err := prepare(url)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf("%s_match_commands", microservice)
	e := NewEmitter()

	err = createHeaderConsumers(q, events)
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

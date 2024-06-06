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
		err := consume(e, q.QueueName)
		if err != nil {
			fmt.Println("Error consuming messages:", err)
		}
	}()

	return e, nil
}

package saga

import (
	"encoding/json"
	"fmt"
	"github.com/jym272/go_library_test/event"
	"slices"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventHandler struct {
	Channel *EventsConsumeChannel  `json:"channel"`
	Payload map[string]interface{} `json:"payload"`
}

func ParseEventPayload[T any](handlerPayload map[string]interface{}, data *T) *T {
	body, err := json.Marshal(handlerPayload)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(body, &data); err != nil {
		panic(err)
	}
	return data
}

// ParseEventPayload (this) It also works, but you need to pass a reference to the variable
// and is not type safe to assure that, as the type is: any
// Works:
// var eventPayload1 saga.SocialNewUserPayload   // or a pointer *saga.SocialNewUserPayload
// ------------------------->key, pass the reference<-----------------//
// handler.ParseEventPayload(&eventPayload1)
//
// It does not work:
// handler.ParseEventPayload(eventPayload1)
func (e *EventHandler) ParseEventPayload(data any) {
	body, err := json.Marshal(e.Payload)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(body, &data); err != nil {
		panic(err)
	}
}

// eventCallback handles the consumption and processing of microservice events.
func eventCallback(msg *amqp.Delivery, channel *amqp.Channel, emitter *Emitter[EventHandler, event.MicroserviceEvent], queueName string) {
	if msg == nil {
		fmt.Println("Message not available")
		return
	}

	var eventPayload map[string]interface{}
	if err := json.Unmarshal(msg.Body, &eventPayload); err != nil {
		fmt.Printf("Error parsing message: %s\n", err)
		err = channel.Nack(msg.DeliveryTag, false, false)
		if err != nil {
			fmt.Println("Error negatively acknowledging message:", err)
			return
		}
		return
	}

	eventKey, err := findEventValues(msg.Headers)
	if err != nil {
		fmt.Println("Invalid header value: no valid event key found")
		err = channel.Nack(msg.DeliveryTag, false, false)
		if err != nil {
			fmt.Println("Error negatively acknowledging message:", err)
			return
		}
		return
	}
	if len(eventKey) != 1 {
		fmt.Println("More then one valid header, using the first one detected, that is because the payload is typed with a particular event")

	}

	responseChannel := &EventsConsumeChannel{
		&ConsumeChannel{
			channel:   channel,
			msg:       msg,
			queueName: queueName,
		},
	}

	emitter.Emit(eventKey[0], EventHandler{
		Payload: eventPayload,
		Channel: responseChannel,
	})
}

// findEventValues find all the MicroserviceEvent values in the headers.
func findEventValues(headers amqp.Table) ([]event.MicroserviceEvent, error) {
	var eventValues []event.MicroserviceEvent
	for _, value := range headers {
		if _, ok := value.(string); !ok {
			continue
		}
		val := event.MicroserviceEvent(value.(string))
		if slices.Contains(event.MicroserviceEventValues(), val) {
			eventValues = append(eventValues, val)
		}
	}
	if len(eventValues) == 0 {
		return nil, fmt.Errorf("no valid event key found")
	}
	return eventValues, nil
}

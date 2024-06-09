package saga

import (
	"encoding/json"
	"fmt"
	"github.com/jym272/go_library_test/micro"
	amqp "github.com/rabbitmq/amqp091-go"
)

type CommandHandler struct {
	Channel *MicroserviceConsumeChannel `json:"channel"`
	Payload map[string]interface{}      `json:"payload"`
	SagaID  int                         `json:"sagaId"`
}

func sagaCommandCallback(msg *amqp.Delivery, channel *amqp.Channel, e *Emitter[CommandHandler, micro.StepCommand], queueName string) {
	if msg == nil {
		fmt.Println("NO MSG AVAILABLE")
		return
	}

	var currentStep SagaStep
	err := json.Unmarshal(msg.Body, &currentStep)
	if err != nil {
		fmt.Println("ERROR PARSING MSG", err)
		err = channel.Nack(msg.DeliveryTag, false, false)
		if err != nil {
			fmt.Println("Error negatively acknowledging message:", err)
			return
		}
		return
	}

	responseChannel := &MicroserviceConsumeChannel{
		step: currentStep,
		ConsumeChannel: &ConsumeChannel{
			channel:   channel,
			msg:       msg,
			queueName: queueName,
		},
	}

	e.Emit(currentStep.Command, CommandHandler{
		Channel: responseChannel,
		Payload: currentStep.PreviousPayload,
		SagaID:  currentStep.SagaID,
	})
}

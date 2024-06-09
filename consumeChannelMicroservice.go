package saga

import (
	"fmt"
)

type (
	StepHashId string
	Occurrence int
)

type MicroserviceConsumeChannel struct {
	*ConsumeChannel
	step SagaStep
}

func (m *MicroserviceConsumeChannel) AckMessage(payloadForNextStep map[string]interface{}) {
	m.step.Status = Success
	previousPayload := m.step.PreviousPayload
	metaData := make(map[string]interface{})

	for key, value := range previousPayload {
		if len(key) > 2 && key[:2] == "__" {
			metaData[key] = value
		}
	}

	for key, value := range payloadForNextStep {
		metaData[key] = value
	}

	m.step.Payload = metaData

	err := sendToQueue(ReplyToSagaQ, m.step)
	if err != nil {
		// TODO: reenqueue message o manejar mejor el error
		return
	}

	err = m.channel.Ack(m.msg.DeliveryTag, false)
	if err != nil {
		// TODO: reenqueue message
		fmt.Println("Error acknowledging message:", err)
	}
}

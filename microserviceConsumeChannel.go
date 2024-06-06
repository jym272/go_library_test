package saga

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	StepHashId string
	Occurrence int
)

var sagaStepOccurrence = make(map[StepHashId]Occurrence)

type MicroserviceConsumeChannel struct {
	channel   *amqp.Channel
	msg       *amqp.Delivery
	queueName string
	step      SagaStep
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

const (
	NACKING_DELAY_MS = 5000 // 5 seconds
	MAX_NACK_RETRIES = 3
	// MAX_OCCURRENCE
	/**
	 * Define the maximum occurrence in a fail saga step of the nack delay with fibonacci strategy
	 * | Occurrence | Delay in the next nack |
	 * |------------|------------------------|
	 * | 17         | 0.44 hours  |
	 * | 18         | 0.72 hours  |
	 * | 19         | 1.18 hours  |
	 * | 20         | 1.88 hours  |
	 * | 21         | 3.04 hours  |
	 * | 22         | 4.92 hours  |
	 * | 23         | 7.96 hours  |
	 * | 24         | 12.87 hours |
	 * | 25         | 20.84 hours |
	 */
	MAX_OCCURRENCE = 19
)

func (m *MicroserviceConsumeChannel) NackWithDelayAndRetries(delay time.Duration, maxRetries int) (int, error) {
	return nackWithDelay(m.msg, m.queueName, delay, maxRetries)
}

func (m *MicroserviceConsumeChannel) NackWithFibonacciStrategy(maxOccurrence int, salt string) (int, int, int, error) {
	// Update saga step occurrence
	hashID := m.getStepHashId(salt)
	occurrence, exists := sagaStepOccurrence[hashID]
	if !exists {
		occurrence = 0
	}
	if int(occurrence) >= maxOccurrence {
		occurrence = 0
	}
	sagaStepOccurrence[hashID] = occurrence + 1

	// Calculate delay
	delay := fibonacci(int(occurrence)) * 1000 // in milliseconds

	// Perform nack with delay and retries
	count, err := m.NackWithDelayAndRetries(time.Duration(delay)*time.Millisecond, math.MaxInt64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error performing nack with fibonacci strategy: %w", err)
	}

	return count, delay, int(occurrence), nil
}

func (m *MicroserviceConsumeChannel) getStepHashId(salt string) StepHashId {
	// Generate hash ID for saga step
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%d-%s-%v-%s", m.step.SagaID, m.step.Command, m.step.Payload, salt)))
	return StepHashId(hex.EncodeToString(hash.Sum(nil))[:10])
}

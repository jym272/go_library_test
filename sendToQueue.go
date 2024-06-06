package saga

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var sendChannel *amqp.Channel

func getSendChannel() (*amqp.Channel, error) {
	if sendChannel == nil {
		conn, err := getRabbitMQConn()
		if err != nil {
			return nil, err
		}
		sendChannel, err = conn.Channel()
		if err != nil {
			return nil, err
		}
	}
	return sendChannel, nil
}

func sendToQueue(queueName Queue, step SagaStep) error {
	err := send(string(queueName), step)
	if err != nil {
		return err
	}
	return nil
}

func send(queueName string, payload interface{}) error {
	channel, err := getSendChannel()
	if err != nil {
		return err
	}

	_, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = channel.PublishWithContext(ctx, "", queueName, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		return err
	}

	return nil
}

func closeSendChannel() error {
	if sendChannel != nil {
		err := sendChannel.Close()
		sendChannel = nil
		return err
	}
	return nil
}

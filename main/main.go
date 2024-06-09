package main

import (
	"fmt"
	"github.com/jym272/go_library_test"
	"github.com/jym272/go_library_test/event"
	"github.com/jym272/go_library_test/micro"
	"time"
)

func main() {

	waitChannel := make(chan struct{})

	transactional := saga.Transactional{
		RabbitUri:    "amqp://rabbit:1234@localhost:5672",
		Microservice: micro.Auth,
		Events: []event.MicroserviceEvent{
			event.SocialNewUserEvent,
			event.SocialBlockChatEvent,
		},
	}
	eventEmitter := transactional.ConnectToEvents()
	eventEmitter.On(event.SocialNewUserEvent, func(handler saga.EventHandler) {
		// el defer tiene que ser obligatorio!!!!
		eventPayload := saga.ParseEventPayload(handler.Payload, &event.SocialNewUserPayload{})
		fmt.Println("SocialNewUserEvent received", eventPayload)
		count, delay, _ := handler.Channel.NackWithDelay(2*time.Second, 4)
		fmt.Println("NackWithDelay", count, delay)
	})
	eventEmitter.On(event.SocialBlockChatEvent, func(handler saga.EventHandler) {
		eventPayload := saga.ParseEventPayload(handler.Payload, &event.SocialBlockChatPayload{})

		fmt.Println("SocialBlockChatEvent received", eventPayload)
		handler.Channel.AckMessage()
	})

	commandEmitter := transactional.ConnectToSagaCommandEmitter()
	commandEmitter.On(micro.MintImageCommand, func(handler saga.CommandHandler) {
		payload := handler.Payload
		fmt.Println("SocialNewUserEvent received", payload)
		handler.Channel.AckMessage(map[string]interface{}{
			"email": "sadfsdf",
		})
	})

	fmt.Println("Transactional started")
	<-waitChannel
}

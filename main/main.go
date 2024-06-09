package main

import (
	"fmt"
	"github.com/jym272/go_library_test"
	"time"
)

func main() {

	waitChannel := make(chan struct{})
	eventEmitter, _, err := saga.StartTransactional(saga.TransactionalConfig{
		Url:          "amqp://rabbit:1234@localhost:5672",
		Microservice: saga.Auth,
		Events: []saga.MicroserviceEvent{
			saga.SocialNewUserEvent,
			saga.SocialBlockChatEvent,
		},
	})
	if err != nil {
		fmt.Println("Error starting transactional:", err)
		return
	}
	eventEmitter.On(saga.SocialNewUserEvent, func(handler saga.EventHandler) {
		// el defer tiene que ser obligatorio!!!!
		eventPayload := saga.ParseEventPayload(handler.Payload, &saga.SocialNewUserPayload{})
		fmt.Println("SocialNewUserEvent received", eventPayload)
		count, delay, _ := handler.Channel.NackWithDelay(2*time.Second, 4)
		fmt.Println("NackWithDelay", count, delay)
	})
	eventEmitter.On(saga.SocialBlockChatEvent, func(handler saga.EventHandler) {
		eventPayload := saga.ParseEventPayload(handler.Payload, &saga.SocialBlockChatPayload{})

		fmt.Println("SocialBlockChatEvent received", eventPayload)
		handler.Channel.AckMessage()
	})

	//commandEmitter.On(saga.NewUserSetRolesToRoomsCommand, func(handler saga.CommandHandler) {
	//	payload := handler.Payload
	//	fmt.Println("SocialNewUserEvent received", payload)
	//	handler.Channel.AckMessage(map[string]interface{}{
	//		"email": "sadfsdf",
	//	})
	//})

	fmt.Println("Transactional started")
	<-waitChannel
}

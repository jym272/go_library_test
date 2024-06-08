package main

import (
	"fmt"
	"github.com/jym272/go_library_test"
)

func masin() {
	err := saga.Prepare("amqp://rabbit:1234@localhost:5672")
	if err != nil {
		fmt.Println("Error preparing:", err)
		return
	}

	err = saga.PublishEvent(saga.SocialNewUserPayload{
		UserID: "123123",
	})
	if err != nil {
		fmt.Println("Error publishing event:", err)
		return
	}
	err = saga.PublishEvent(saga.SocialBlockChatPayload{
		UserID:        "11d1",
		UserToBlockID: "1d2d12d",
	})
	if err != nil {
		fmt.Println("Error publishing event:", err)
		return
	}

}

func main() {

	waitChannel := make(chan struct{})
	eventEmitter, commandEmitter, err := saga.StartTransactional(saga.TransactionalConfig{
		Url:          "amqp://rabbit:1234@localhost:5672",
		Microservice: saga.RoomCreator,
		Events: []saga.MicroserviceEvent{
			saga.SocialNewUserEvent,
			saga.SocialBlockChatEvent,
		},
	})
	if err != nil {
		fmt.Println("Error starting transactional:", err)
		return
	}
	// se puede agregar validación en runtime a aSocialNewUserEvent
	eventEmitter.On(saga.SocialNewUserEvent, func(handler saga.EventHandler) {
		eventPayload := saga.ParseEventPayload(handler.Payload, &saga.SocialNewUserPayload{})
		//var eventPayload1 saga.SocialNewUserPayload
		//handler.ParseEventPayload(&eventPayload1)
		//fmt.Println("SocialNewUserEvent received", eventPayload1)
		fmt.Println("SocialNewUserEvent received", eventPayload)
		handler.Channel.AckMessage()
	})
	eventEmitter.On(saga.SocialBlockChatEvent, func(handler saga.EventHandler) {
		eventPayload := saga.ParseEventPayload(handler.Payload, &saga.SocialBlockChatPayload{})

		fmt.Println("SocialBlockChatEvent received", eventPayload)
		handler.Channel.AckMessage()
	})

	commandEmitter.On(saga.NewUserSetRolesToRoomsCommand, func(handler saga.CommandHandler) {
		payload := handler.Payload
		fmt.Println("SocialNewUserEvent received", payload)
	})

	fmt.Println("Transactional started")
	<-waitChannel
}

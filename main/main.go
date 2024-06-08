package main

import (
	"fmt"
	"github.com/jym272/go_library_test"
)

func main() {
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

}

func mai1n() {

	waitChannel := make(chan struct{})
	eventEmitter, commandEmitter, err := saga.StartTransactional(saga.TransactionalConfig{
		Url:          "amqp://rabbit:1234@localhost:5672",
		Microservice: saga.RoomCreator,
		Events: []saga.MicroserviceEvent{
			saga.SocialNewUserEvent,
		},
	})
	if err != nil {
		fmt.Println("Error starting transactional:", err)
		return
	}
	eventEmitter.On(saga.SocialNewUserEvent, func(handler saga.EventHandler) {
		payload := handler.Payload.(saga.SocialNewUserPayload)
		fmt.Println("SocialNewUserEvent received", payload.UserID)
		handler.Channel.AckMessage()
	})

	commandEmitter.On(saga.NewUserSetRolesToRoomsCommand, func(handler saga.CommandHandler) {
		payload := handler.Payload
		fmt.Println("SocialNewUserEvent received", payload)
	})

	fmt.Println("Transactional started")
	<-waitChannel
}

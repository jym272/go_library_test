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
	err = saga.PublishEvent(saga.SocialBlockChatPayload{
		UserID:        "11d1",
		UserToBlockID: "1d2d12d",
	})
	if err != nil {
		fmt.Println("Error publishing event:", err)
		return
	}

	//err = saga.CommenceSaga(saga.UpdateUserImagePayload{
	//	UserId:     "1d2",
	//	FolderName: "12d",
	//	BucketName: "12d",
	//})
	//if err != nil {
	//	fmt.Println("Error commencing saga:", err)
	//	return
	//}

}

func masdin() {

	waitChannel := make(chan struct{})
	eventEmitter, commandEmitter, err := saga.StartTransactional(saga.TransactionalConfig{
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
		handler.Channel.AckMessage(map[string]interface{}{
			"email": "sadfsdf",
		})
	})

	fmt.Println("Transactional started")
	<-waitChannel
}

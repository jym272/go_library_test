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

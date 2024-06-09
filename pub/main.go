package main

import (
	"crypto/rand"
	"fmt"
	"github.com/jym272/go_library_test"
	"github.com/jym272/go_library_test/event"
	"math/big"
)

func GenerateRandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return ""
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret)
}

func main() {

	saga.Prepare("amqp://rabbit:1234@localhost:5672")

	i := 0

	for i < 10 {

		err := saga.PublishEvent(event.SocialNewUserPayload{
			UserID: GenerateRandomString(5),
		})
		if err != nil {
			fmt.Println("Error publishing event:", err)
			return
		}

		err = saga.PublishEvent(event.SocialBlockChatPayload{
			UserID:        GenerateRandomString(5),
			UserToBlockID: "1d2d12d",
		})
		if err != nil {
			fmt.Println("Error publishing event:", err)
			return
		}

		i++
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

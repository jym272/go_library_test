package saga

type MicroserviceEvent string

const (
	TestImage            MicroserviceEvent = "test.image"
	TestMint             MicroserviceEvent = "test.mint"
	PaymentsNotifyClient MicroserviceEvent = "payments.notify_client"
	SocialBlockChat      MicroserviceEvent = "social.block_chat"
	SocialNewUser        MicroserviceEvent = "social.new_user"
	SocialUnblockChat    MicroserviceEvent = "social.unblock_chat"
)

func microserviceEventValues() []MicroserviceEvent {
	return []MicroserviceEvent{
		TestImage,
		TestMint,
		PaymentsNotifyClient,
		SocialBlockChat,
		SocialNewUser,
		SocialUnblockChat,
	}
}

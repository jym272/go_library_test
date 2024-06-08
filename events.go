package saga

type MicroserviceEvent string

const (
	TestImageEvent MicroserviceEvent = "test.image"
	TestMintEvent  MicroserviceEvent = "test.mint"

	PaymentsNotifyClientEvent MicroserviceEvent = "payments.notify_client"
	SocialBlockChatEvent      MicroserviceEvent = "social.block_chat"
	SocialNewUserEvent        MicroserviceEvent = "social.new_user"
	SocialUnblockChatEvent    MicroserviceEvent = "social.unblock_chat"
)

type PayloadEvent interface {
	Type() MicroserviceEvent // Discriminant method
}

type SocialUnblockChatPayload struct {
	UserID          string `json:"userId"`
	UserToUnblockID string `json:"userToUnblockId"`
}

func (SocialUnblockChatPayload) Type() MicroserviceEvent {
	return SocialUnblockChatEvent
}

type PaymentsNotifyClientPayload struct {
	Room    string                 `json:"room"`
	Message map[string]interface{} `json:"message"`
}

func (PaymentsNotifyClientPayload) Type() MicroserviceEvent {
	return PaymentsNotifyClientEvent
}

//	func getEvent(payload PayloadEvent) MicroserviceEvent {
//		return payload.Type()
//		//case PaymentsNotifyClientEvent:
//		//	textPayload := payload.(PaymentsNotifyClientPayload)
//		//	// Handle text payload
//		//case SocialUnblockChatEvent:
//		//	imagePayload := payload.(ImagePayload)
//		//	// Handle image payload
//		//default:
//		//	// Handle unknown payload types
//		//}
//	}
func microserviceEventValues() []MicroserviceEvent {
	return []MicroserviceEvent{
		TestImageEvent,
		TestMintEvent,

		PaymentsNotifyClientEvent,
		SocialBlockChatEvent,
		SocialNewUserEvent,
		SocialUnblockChatEvent,
	}
}

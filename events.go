package saga

type MicroserviceEvent string
type PayloadEvent interface {
	Type() MicroserviceEvent // Discriminant method
}

const (
	TestImageEvent MicroserviceEvent = "test.image"
	TestMintEvent  MicroserviceEvent = "test.mint"

	PaymentsNotifyClientEvent MicroserviceEvent = "payments.notify_client"
	SocialBlockChatEvent      MicroserviceEvent = "social.block_chat"
	SocialNewUserEvent        MicroserviceEvent = "social.new_user"
	SocialUnblockChatEvent    MicroserviceEvent = "social.unblock_chat"
)

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

type AllPayloads interface {
	TestImagePayload | TestMintPayload | SocialBlockChatPayload | SocialNewUserPayload | SocialUnblockChatPayload | PaymentsNotifyClientPayload
}

type TestImagePayload struct {
	Image string `json:"image"`
}

func (TestImagePayload) Type() MicroserviceEvent {
	return TestImageEvent
}

type TestMintPayload struct {
	Mint string `json:"mint"`
}

func (TestMintPayload) Type() MicroserviceEvent {
	return TestMintEvent
}

type SocialBlockChatPayload struct {
	UserID        string `json:"userId"`
	UserToBlockID string `json:"userToBlockId"`
}

func (SocialBlockChatPayload) Type() MicroserviceEvent {
	return SocialBlockChatEvent
}

type SocialNewUserPayload struct {
	UserID string `json:"userId"`
}

func (SocialNewUserPayload) Type() MicroserviceEvent {
	return SocialNewUserEvent
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

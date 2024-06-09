package saga

type AvailableMicroservices string

const (
	Auth        AvailableMicroservices = "auth"
	RoomCreator AvailableMicroservices = "room-creator"
)

type (
	Exchange string
	Queue    string
)

const (
	RequeueExchange         Exchange = "requeue_exchange"
	CommandsExchange        Exchange = "commands_exchange"
	MatchingExchange        Exchange = "matching_exchange"
	MatchingRequeueExchange Exchange = "matching_requeue_exchange"
)

const (
	ReplyToSagaQ Queue = "reply_to_saga"
)

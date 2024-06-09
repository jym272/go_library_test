package saga

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

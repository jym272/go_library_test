package micro

type StepCommand = string

type AvailableMicroservices string

const (
	Auth           AvailableMicroservices = "auth"
	RoomCreator    AvailableMicroservices = "room-creator"
	Showcase       AvailableMicroservices = "legend-showcase"
	RapidMessaging AvailableMicroservices = "rapid-messaging"
)

// image
const (
	TestImage          AvailableMicroservices = "test-image"
	CreateImageCommand StepCommand            = "create_image"
	UpdateTokenCommand StepCommand            = "update_token"
)

// mint
const (
	TestMint         AvailableMicroservices = "test-mint"
	MintImageCommand StepCommand            = "mint_image"
)

// payments
const (
	Payments                               AvailableMicroservices = "payments"
	ResourcePurchasedDeductCoinsCommand    StepCommand            = "resource_purchased:deduct_coins"
	ResourcePurchasedRemoveRedisKeyCommand StepCommand            = "resource_purchased:remove_redis_key"
)

// room-inventory
const (
	RoomInventory                    AvailableMicroservices = "room-inventory"
	DecreaseAvailableQuantityCommand StepCommand            = "resource_purchased:decrease_available_quantity"
)

// room-snapshot
const (
	RoomSnapshot            AvailableMicroservices = "room-snapshot"
	PurchaseResourceCommand StepCommand            = "resource_purchased:save_purchased_resource"
)

// social
const (
	Social                 AvailableMicroservices = "social"
	UpdateUserImageCommand StepCommand            = "update_user:image"
	NotifyClientCommand    StepCommand            = "notify_client"
)

// storage
const (
	Storage           AvailableMicroservices = "legend-storage"
	UpdateFileCommand StepCommand            = "update_file"
)

package micro

type StepCommand = string

type AvailableMicroservices string

const (
	Auth           AvailableMicroservices = "auth"
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
	Payments                            AvailableMicroservices = "payments"
	ResourcePurchasedDeductCoinsCommand StepCommand            = "resource_purchased:deduct_coins"
)

// room-inventory
const (
	RoomInventory                AvailableMicroservices = "room-inventory"
	SavePurchasedResourceCommand StepCommand            = "resource_purchased:save_purchased_resource"
)

// room-snapshot
const (
	RoomSnapshot AvailableMicroservices = "room-snapshot"
)

// room-creator
const (
	RoomCreator AvailableMicroservices = "room-creator"
)

// legend-showcase
const (
	Showcase AvailableMicroservices = "legend-showcase"
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

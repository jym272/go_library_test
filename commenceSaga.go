package saga

const (
	CommenceSagaQueue Queue = "commence_saga"
)

type SagaTitle string

const (
	UpdateUserImage      SagaTitle = "update_user:image"
	PurchaseResourceFlow SagaTitle = "purchase_resource_flow"
)

type CommencePayload interface {
	Type() SagaTitle
}

// UpdateUserImagePayload is the payload for the update_user:image event.
type UpdateUserImagePayload struct {
	UserId     string `json:"userId"`
	FolderName string `json:"folderName"`
	BucketName string `json:"bucketName"`
}

func (UpdateUserImagePayload) Type() SagaTitle {
	return UpdateUserImage
}

// PurchaseResourceFlowPayload is the payload for the purchase_resource_flow event.
type PurchaseResourceFlowPayload struct {
	UserId     string `json:"userId"`
	ResourceId string `json:"resourceId"`
	Price      int    `json:"price"`
	Quantity   int    `json:"quantity"`
}

func (PurchaseResourceFlowPayload) Type() SagaTitle {
	return PurchaseResourceFlow
}

type commenceSaga struct {
	Title   SagaTitle   `json:"title"`
	Payload interface{} `json:"payload"`
}

func CommenceSaga(payload CommencePayload) error {
	title := payload.Type()
	err := send(string(CommenceSagaQueue), commenceSaga{
		Title:   title,
		Payload: payload,
	})
	if err != nil {
		return err
	}
	return nil
}

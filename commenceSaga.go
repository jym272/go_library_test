package saga

const (
	CommenceSagaQueue Queue = "commence_saga"
)

type (
	SagaTitle string
)

const (
	RechargeBalance SagaTitle = "recharge_balance"
	BuyProducts     SagaTitle = "buy_products"
	PersistRedsys   SagaTitle = "persist_redsys"
	SendEmailToken  SagaTitle = "send_email:token"
)

type commenceSaga struct {
	Title   SagaTitle   `json:"title"`
	Payload interface{} `json:"payload"`
}

func CommenceSaga(title SagaTitle, payload interface{}) error {
	err := send(string(CommenceSagaQueue), commenceSaga{
		Title:   title,
		Payload: payload,
	})
	if err != nil {
		return err
	}
	return nil
}

package dto

type AddTransactionPayload struct {
	OrderId int `json:"order_id" valid:"required~Order id can't be empty"`
}

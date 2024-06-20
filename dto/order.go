package dto

type AddOrderPayload struct {
	ProductId int `json:"product_id" valid:"required~Product id can't be empty"`
	Qty       int `json:"qty" valid:"required~Qty can't be empty"`
}

type ModifyOrderPayload struct {
	Qty int `json:"qty" valid:"required~Qty can't be empty"`
}

package dto

type CategoryPayload struct {
	Type string `json:"type" valid:"required~Type can't be empty"`
}

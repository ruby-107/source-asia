package model

type Request struct {
	UserID  int    `json:"user_id" validate:"required,gt=0"`
	Payload string `json:"payload"`
}

package models

type Token struct {
	Code string `json:"code" form:"code" validate:"required"`
}

package models

import "github.com/go-playground/validator/v10"

type Position struct {
	PositionId uint   `gorm:"primary_key"`
	Name       string `validate:"required"`
}

func (u *Position) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

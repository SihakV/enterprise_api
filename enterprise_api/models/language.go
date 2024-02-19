package models

import "github.com/go-playground/validator/v10"

type Language struct {
	LanguageId uint   `gorm:"primary_key"`
	Name       string `validate:"required"`
}

func (u *Language) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

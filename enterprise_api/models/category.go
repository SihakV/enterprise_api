package models

import "github.com/go-playground/validator/v10"

type Category struct {
	CategoryId uint   `gorm:"primary_key"`
	Name       string `validate:"required"`
}

func (Category) TableName() string {
	return "categories"
}

func (u *Category) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

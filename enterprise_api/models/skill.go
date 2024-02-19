package models

import "github.com/go-playground/validator/v10"

type Skill struct {
	SkillId uint   `gorm:"primary_key"`
	Name    string `validate:"required"`
}

func (u *Skill) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

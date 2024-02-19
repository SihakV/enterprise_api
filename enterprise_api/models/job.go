package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Job struct {
	JobId             uint   `gorm:"primary_key"`
	Title             string `validate:"required"`
	CategoryIds       string `validate:"required"`
	Description       string `validate:"required"`
	Contact           string `validate:"required"`
	ExpiryDate        string `validate:"required"`
	AnnouncementImage *string
	CreatedBy         int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Status            int
}

func (u *Job) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

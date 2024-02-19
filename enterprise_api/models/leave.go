package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Leave struct {
	LeaveId       uint   `gorm:"primary_key"`
	Title         string `validate:"required"`
	LeaveFile     *string
	Description   string `validate:"required"`
	LeaveDateFrom string `validate:"required"`
	LeaveDateTo   string `validate:"required"`
	CreatedBy     int
	ApprovedBy    int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Status        int
}

func (Leave) TableName() string {
	return "leave_requests"
}

func (u *Leave) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

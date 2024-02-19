package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Expense struct {
	ExpenseId   uint    ` gorm:"primary_key"`
	Title       string  `validate:"required"`
	Amount      float64 `validate:"required"`
	ExpenseFile *string
	Description string
	CreatedBy   int
	ApprovedBy  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      int
}

func (Expense) TableName() string {
	return "expense_requests"
}

func (u *Expense) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

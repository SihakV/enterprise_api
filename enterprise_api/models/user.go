package models

import (
	"errors"
	"midterm/db"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
)

type User struct {
	UserId     uint   `gorm:"primary_key"`
	Username   string `validate:"required"`
	Email      string `validate:"required,email"`
	RoleId     int    `validate:"required,max=2"`
	Dob        string
	PositionId int `validate:"required"`
	Phone      string
	Salary     float32 `validate:"required,numeric"`
	Profile    *string
	Password   string `validate:"required,min=5"`
	CreatedBy  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Status     int
}

func (u *User) Validate() error {
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		return err
	}

	// Check for unique email
	var result User
	if err := db.DbConnect.Where("email = ?", u.Email).First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	} else {
		return errors.New("email already in use")
	}

	return nil
}

package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Applicant struct {
	ApplicantId   uint   `gorm:"primary_key"`
	Name          string `validate:"required"`
	Email         string `validate:"required"`
	Phone         string `validate:"required"`
	ScheduledDate time.Time
	LanguageIds   *string
	SkillIds      *string
	Address       *string
	Experience    *string
	Education     *string
	Summary       *string
	JobId         int
	CvFile        *string
	ApprovedBy    int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Status        int
}

func (u *Applicant) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

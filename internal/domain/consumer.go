package domain

import (
	"errors"
	"regexp"
	"time"
)

var (
	ErrInvalidNIK       = errors.New("invalid NIK: must be 16 digits")
	ErrInvalidFullName  = errors.New("invalid full name")
	ErrInvalidBirthDate = errors.New("invalid birth date")
	ErrInvalidSalary    = errors.New("invalid salary: must be positive")
	ErrConsumerNotFound = errors.New("consumer not found")
)

// Consumer represents a customer in the system
type Consumer struct {
	ID          int64     `json:"id"`
	NIK         string    `json:"nik"`
	FullName    string    `json:"full_name"`
	LegalName   string    `json:"legal_name"`
	BirthPlace  string    `json:"birth_place"`
	BirthDate   time.Time `json:"birth_date"`
	Salary      float64   `json:"salary"`
	KTPPhoto    string    `json:"ktp_photo"`
	SelfiePhoto string    `json:"selfie_photo"`
	CreatedAt   time.Time `json:"created_at"`
}

// Validate performs validation on Consumer data
func (c *Consumer) Validate() error {
	// Validate NIK (must be exactly 16 digits)
	nikPattern := regexp.MustCompile(`^\d{16}$`)
	if !nikPattern.MatchString(c.NIK) {
		return ErrInvalidNIK
	}

	// Validate full name (not empty, max 100 chars)
	if len(c.FullName) == 0 || len(c.FullName) > 100 {
		return ErrInvalidFullName
	}

	// Validate salary (must be positive)
	if c.Salary <= 0 {
		return ErrInvalidSalary
	}

	// Validate birth date (must be in the past)
	if c.BirthDate.After(time.Now()) {
		return ErrInvalidBirthDate
	}

	return nil
}

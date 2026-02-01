package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidTenorValue = errors.New("invalid tenor: must be 1, 2, 3, or 6")
	ErrLimitNotFound     = errors.New("consumer limit not found")
	ErrInsufficientLimit = errors.New("insufficient limit available")
)

// ConsumerLimit represents the credit limit for a consumer at a specific tenor
type ConsumerLimit struct {
	ID          int64     `json:"id"`
	ConsumerID  int64     `json:"consumer_id"`
	Tenor       int       `json:"tenor"`
	LimitAmount float64   `json:"limit_amount"`
	UsedAmount  float64   `json:"used_amount"`
	Version     int       `json:"version"` // For optimistic locking
	UpdatedAt   time.Time `json:"updated_at"`
}

// AvailableLimit returns the remaining available limit
func (cl *ConsumerLimit) AvailableLimit() float64 {
	return cl.LimitAmount - cl.UsedAmount
}

// CanAccommodate checks if the limit can accommodate the requested amount
func (cl *ConsumerLimit) CanAccommodate(amount float64) bool {
	return cl.AvailableLimit() >= amount
}

// ValidateTenor checks if the tenor is valid
func ValidateTenor(tenor int) error {
	validTenors := map[int]bool{1: true, 2: true, 3: true, 6: true}
	if !validTenors[tenor] {
		return ErrInvalidTenorValue
	}
	return nil
}

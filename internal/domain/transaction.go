package domain

import (
	"errors"
	"time"
)

var (
	ErrLimitExceeded       = errors.New("consumer limit exceeded")
	ErrInvalidTenor        = errors.New("invalid tenor")
	ErrInvalidOTR          = errors.New("invalid OTR: must be positive")
	ErrInvalidAdminFee     = errors.New("invalid admin fee: must be non-negative")
	ErrInvalidInstallment  = errors.New("invalid installment amount: must be positive")
	ErrInvalidAssetName    = errors.New("invalid asset name")
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrDuplicateContract   = errors.New("duplicate contract number")
)

// Transaction represents a financing transaction
type Transaction struct {
	ID                int64     `json:"id"`
	ContractNumber    string    `json:"contract_number"`
	ConsumerID        int64     `json:"consumer_id"`
	Tenor             int       `json:"tenor"`
	OTR               float64   `json:"otr"` // On The Road price
	AdminFee          float64   `json:"admin_fee"`
	InstallmentAmount float64   `json:"installment_amount"`
	InterestAmount    float64   `json:"interest_amount"`
	AssetName         string    `json:"asset_name"`
	CreatedAt         time.Time `json:"created_at"`
}

// Validate performs validation on Transaction data
func (t *Transaction) Validate() error {
	// Validate tenor
	if err := ValidateTenor(t.Tenor); err != nil {
		return err
	}

	// Validate OTR (must be positive)
	if t.OTR <= 0 {
		return ErrInvalidOTR
	}

	// Validate admin fee (must be non-negative)
	if t.AdminFee < 0 {
		return ErrInvalidAdminFee
	}

	// Validate installment amount (must be positive)
	if t.InstallmentAmount <= 0 {
		return ErrInvalidInstallment
	}

	// Validate asset name (not empty, max 100 chars)
	if len(t.AssetName) == 0 || len(t.AssetName) > 100 {
		return ErrInvalidAssetName
	}

	return nil
}

// TotalAmount returns the total financed amount (OTR + AdminFee + Interest)
func (t *Transaction) TotalAmount() float64 {
	return t.OTR + t.AdminFee + t.InterestAmount
}

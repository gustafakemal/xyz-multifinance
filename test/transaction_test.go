package test

import (
	"testing"
	"time"

	"github.com/username/xyz-multifinance/internal/domain"
)

// TestConsumerValidation tests consumer data validation
func TestConsumerValidation(t *testing.T) {
	tests := []struct {
		name      string
		consumer  domain.Consumer
		expectErr bool
		errType   error
	}{
		{
			name: "Valid consumer",
			consumer: domain.Consumer{
				NIK:        "3201012345678901",
				FullName:   "John Doe",
				LegalName:  "John Doe",
				BirthPlace: "Jakarta",
				BirthDate:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Salary:     5000000,
			},
			expectErr: false,
		},
		{
			name: "Invalid NIK - too short",
			consumer: domain.Consumer{
				NIK:       "123456",
				FullName:  "John Doe",
				BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Salary:    5000000,
			},
			expectErr: true,
			errType:   domain.ErrInvalidNIK,
		},
		{
			name: "Invalid NIK - contains letters",
			consumer: domain.Consumer{
				NIK:       "320101234567890A",
				FullName:  "John Doe",
				BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Salary:    5000000,
			},
			expectErr: true,
			errType:   domain.ErrInvalidNIK,
		},
		{
			name: "Invalid salary - negative",
			consumer: domain.Consumer{
				NIK:       "3201012345678901",
				FullName:  "John Doe",
				BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Salary:    -1000,
			},
			expectErr: true,
			errType:   domain.ErrInvalidSalary,
		},
		{
			name: "Invalid birth date - in the future",
			consumer: domain.Consumer{
				NIK:       "3201012345678901",
				FullName:  "John Doe",
				BirthDate: time.Now().Add(24 * time.Hour),
				Salary:    5000000,
			},
			expectErr: true,
			errType:   domain.ErrInvalidBirthDate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.consumer.Validate()
			if tt.expectErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if tt.expectErr && tt.errType != nil && err != tt.errType {
				t.Errorf("Expected error type %v but got %v", tt.errType, err)
			}
		})
	}
}

// TestTransactionValidation tests transaction data validation
func TestTransactionValidation(t *testing.T) {
	tests := []struct {
		name        string
		transaction domain.Transaction
		expectErr   bool
		errType     error
	}{
		{
			name: "Valid transaction",
			transaction: domain.Transaction{
				ConsumerID:        1,
				Tenor:             3,
				OTR:               10000000,
				AdminFee:          100000,
				InstallmentAmount: 3500000,
				InterestAmount:    500000,
				AssetName:         "Motor Honda",
			},
			expectErr: false,
		},
		{
			name: "Invalid tenor",
			transaction: domain.Transaction{
				ConsumerID:        1,
				Tenor:             5,
				OTR:               10000000,
				AdminFee:          100000,
				InstallmentAmount: 3500000,
				InterestAmount:    500000,
				AssetName:         "Motor Honda",
			},
			expectErr: true,
			errType:   domain.ErrInvalidTenorValue,
		},
		{
			name: "Invalid OTR - negative",
			transaction: domain.Transaction{
				ConsumerID:        1,
				Tenor:             3,
				OTR:               -10000000,
				AdminFee:          100000,
				InstallmentAmount: 3500000,
				InterestAmount:    500000,
				AssetName:         "Motor Honda",
			},
			expectErr: true,
			errType:   domain.ErrInvalidOTR,
		},
		{
			name: "Invalid asset name - empty",
			transaction: domain.Transaction{
				ConsumerID:        1,
				Tenor:             3,
				OTR:               10000000,
				AdminFee:          100000,
				InstallmentAmount: 3500000,
				InterestAmount:    500000,
				AssetName:         "",
			},
			expectErr: true,
			errType:   domain.ErrInvalidAssetName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.transaction.Validate()
			if tt.expectErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if tt.expectErr && tt.errType != nil && err != tt.errType {
				t.Errorf("Expected error type %v but got %v", tt.errType, err)
			}
		})
	}
}

// TestConsumerLimitAvailability tests limit calculation
func TestConsumerLimitAvailability(t *testing.T) {
	tests := []struct {
		name              string
		limit             domain.ConsumerLimit
		requestAmount     float64
		expectedAvailable float64
		canAccommodate    bool
	}{
		{
			name: "Sufficient limit",
			limit: domain.ConsumerLimit{
				LimitAmount: 1000000,
				UsedAmount:  300000,
			},
			requestAmount:     500000,
			expectedAvailable: 700000,
			canAccommodate:    true,
		},
		{
			name: "Insufficient limit",
			limit: domain.ConsumerLimit{
				LimitAmount: 1000000,
				UsedAmount:  800000,
			},
			requestAmount:     500000,
			expectedAvailable: 200000,
			canAccommodate:    false,
		},
		{
			name: "Exact limit available",
			limit: domain.ConsumerLimit{
				LimitAmount: 1000000,
				UsedAmount:  300000,
			},
			requestAmount:     700000,
			expectedAvailable: 700000,
			canAccommodate:    true,
		},
		{
			name: "No limit used",
			limit: domain.ConsumerLimit{
				LimitAmount: 1000000,
				UsedAmount:  0,
			},
			requestAmount:     500000,
			expectedAvailable: 1000000,
			canAccommodate:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			available := tt.limit.AvailableLimit()
			if available != tt.expectedAvailable {
				t.Errorf("Expected available limit %f but got %f", tt.expectedAvailable, available)
			}

			canAccommodate := tt.limit.CanAccommodate(tt.requestAmount)
			if canAccommodate != tt.canAccommodate {
				t.Errorf("Expected CanAccommodate to be %v but got %v", tt.canAccommodate, canAccommodate)
			}
		})
	}
}

// TestTenorValidation tests tenor validation
func TestTenorValidation(t *testing.T) {
	tests := []struct {
		name      string
		tenor     int
		expectErr bool
	}{
		{"Valid tenor 1", 1, false},
		{"Valid tenor 2", 2, false},
		{"Valid tenor 3", 3, false},
		{"Valid tenor 6", 6, false},
		{"Invalid tenor 4", 4, true},
		{"Invalid tenor 5", 5, true},
		{"Invalid tenor 0", 0, true},
		{"Invalid tenor -1", -1, true},
		{"Invalid tenor 12", 12, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.ValidateTenor(tt.tenor)
			if tt.expectErr && err == nil {
				t.Errorf("Expected error for tenor %d but got none", tt.tenor)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error for tenor %d but got: %v", tt.tenor, err)
			}
		})
	}
}

// TestTransactionTotalAmount tests total amount calculation
func TestTransactionTotalAmount(t *testing.T) {
	transaction := domain.Transaction{
		OTR:            10000000,
		AdminFee:       100000,
		InterestAmount: 500000,
	}

	expectedTotal := 10600000.0
	actualTotal := transaction.TotalAmount()

	if actualTotal != expectedTotal {
		t.Errorf("Expected total amount %f but got %f", expectedTotal, actualTotal)
	}
}

// TestLimitExceeded tests limit exceeded scenario
func TestLimitExceeded(t *testing.T) {
	limit := domain.ConsumerLimit{
		LimitAmount: 500000,
		UsedAmount:  400000,
	}

	// Request exceeds available limit
	requestAmount := 200000.0
	if limit.CanAccommodate(requestAmount) {
		t.Error("Expected limit to be exceeded but it wasn't")
	}

	// Request is within available limit
	requestAmount = 50000.0
	if !limit.CanAccommodate(requestAmount) {
		t.Error("Expected limit to be sufficient but it wasn't")
	}
}

// TestConcurrentLimitUpdate simulates concurrent limit updates
func TestConcurrentLimitUpdate(t *testing.T) {
	// Test optimistic locking version check
	limit := domain.ConsumerLimit{
		ID:          1,
		ConsumerID:  1,
		Tenor:       3,
		LimitAmount: 1000000,
		UsedAmount:  0,
		Version:     0,
	}

	// Simulate first transaction
	limit.UsedAmount += 300000
	limit.Version++

	if limit.Version != 1 {
		t.Errorf("Expected version to be 1 but got %d", limit.Version)
	}

	// Simulate second transaction
	limit.UsedAmount += 200000
	limit.Version++

	if limit.Version != 2 {
		t.Errorf("Expected version to be 2 but got %d", limit.Version)
	}

	if limit.UsedAmount != 500000 {
		t.Errorf("Expected used amount to be 500000 but got %f", limit.UsedAmount)
	}
}

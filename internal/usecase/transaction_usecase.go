package usecase

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/username/xyz-multifinance/internal/domain"
	"github.com/username/xyz-multifinance/internal/repository"
)

// TransactionUsecase handles business logic for transactions
type TransactionUsecase struct {
	transactionRepo repository.TransactionRepository
	limitRepo       repository.LimitRepository
	consumerRepo    repository.ConsumerRepository
	db              *sql.DB
}

// NewTransactionUsecase creates a new transaction use case
func NewTransactionUsecase(
	transactionRepo repository.TransactionRepository,
	limitRepo repository.LimitRepository,
	consumerRepo repository.ConsumerRepository,
	db *sql.DB,
) *TransactionUsecase {
	return &TransactionUsecase{
		transactionRepo: transactionRepo,
		limitRepo:       limitRepo,
		consumerRepo:    consumerRepo,
		db:              db,
	}
}

// CreateTransactionRequest represents the request to create a transaction
type CreateTransactionRequest struct {
	ConsumerID        int64   `json:"consumer_id"`
	Tenor             int     `json:"tenor"`
	OTR               float64 `json:"otr"`
	AdminFee          float64 `json:"admin_fee"`
	InstallmentAmount float64 `json:"installment_amount"`
	InterestAmount    float64 `json:"interest_amount"`
	AssetName         string  `json:"asset_name"`
}

// CreateTransaction creates a new transaction with concurrent transaction handling
// This method uses database transactions with FOR UPDATE locks and optimistic locking
// to handle concurrent requests safely (ACID compliance)
func (uc *TransactionUsecase) CreateTransaction(req CreateTransactionRequest) (*domain.Transaction, error) {
	// Validate consumer exists
	_, err := uc.consumerRepo.FindByID(req.ConsumerID)
	if err != nil {
		return nil, fmt.Errorf("invalid consumer: %w", err)
	}

	// Create transaction object
	transaction := &domain.Transaction{
		ContractNumber:    uc.generateContractNumber(),
		ConsumerID:        req.ConsumerID,
		Tenor:             req.Tenor,
		OTR:               req.OTR,
		AdminFee:          req.AdminFee,
		InstallmentAmount: req.InstallmentAmount,
		InterestAmount:    req.InterestAmount,
		AssetName:         req.AssetName,
		CreatedAt:         time.Now(),
	}

	// Validate transaction data
	if err := transaction.Validate(); err != nil {
		return nil, err
	}

	// Begin database transaction for ACID compliance
	tx, err := uc.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Get consumer limit with FOR UPDATE lock (prevents concurrent modifications)
	var limit domain.ConsumerLimit
	query := `
		SELECT id, consumer_id, tenor, limit_amount, used_amount, version, updated_at
		FROM consumer_limits
		WHERE consumer_id = ? AND tenor = ?
		FOR UPDATE
	`
	err = tx.QueryRow(query, req.ConsumerID, req.Tenor).Scan(
		&limit.ID,
		&limit.ConsumerID,
		&limit.Tenor,
		&limit.LimitAmount,
		&limit.UsedAmount,
		&limit.Version,
		&limit.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrLimitNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get consumer limit: %w", err)
	}

	// Calculate total financed amount
	totalAmount := transaction.TotalAmount()

	// Check if consumer has sufficient limit
	if !limit.CanAccommodate(totalAmount) {
		return nil, domain.ErrLimitExceeded
	}

	// Update used amount with optimistic locking
	newUsedAmount := limit.UsedAmount + totalAmount
	err = uc.limitRepo.UpdateUsedAmountWithTx(tx, limit.ID, newUsedAmount, limit.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to update limit: %w", err)
	}

	// Create transaction record
	err = uc.transactionRepo.CreateWithTx(tx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transaction, nil
}

// GetTransactionByContract retrieves a transaction by contract number
func (uc *TransactionUsecase) GetTransactionByContract(contractNumber string) (*domain.Transaction, error) {
	return uc.transactionRepo.FindByContractNumber(contractNumber)
}

// GetConsumerTransactions retrieves all transactions for a consumer
func (uc *TransactionUsecase) GetConsumerTransactions(consumerID int64) ([]*domain.Transaction, error) {
	// Verify consumer exists
	_, err := uc.consumerRepo.FindByID(consumerID)
	if err != nil {
		return nil, err
	}

	return uc.transactionRepo.FindByConsumerID(consumerID)
}

// GetConsumerLimits retrieves all limits for a consumer
func (uc *TransactionUsecase) GetConsumerLimits(consumerID int64) ([]*domain.ConsumerLimit, error) {
	// Verify consumer exists
	_, err := uc.consumerRepo.FindByID(consumerID)
	if err != nil {
		return nil, err
	}

	return uc.limitRepo.FindAllByConsumerID(consumerID)
}

// generateContractNumber generates a unique contract number
func (uc *TransactionUsecase) generateContractNumber() string {
	// Format: TRX-YYYYMMDD-UUID
	dateStr := time.Now().Format("20060102")
	uuid := uuid.New().String()[:8]
	return fmt.Sprintf("TRX-%s-%s", dateStr, uuid)
}

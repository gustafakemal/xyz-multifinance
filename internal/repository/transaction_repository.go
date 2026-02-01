package repository

import (
	"database/sql"
	"fmt"

	"github.com/username/xyz-multifinance/internal/domain"
)

// TransactionRepository defines methods for transaction data access
type TransactionRepository interface {
	Create(tx *domain.Transaction) error
	FindByContractNumber(contractNumber string) (*domain.Transaction, error)
	FindByConsumerID(consumerID int64) ([]*domain.Transaction, error)
	BeginTx() (*sql.Tx, error)
	CreateWithTx(sqlTx *sql.Tx, tx *domain.Transaction) error
}

type transactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(tx *domain.Transaction) error {
	query := `
		INSERT INTO transactions (contract_number, consumer_id, tenor, otr, admin_fee,
		                          installment_amount, interest_amount, asset_name)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		tx.ContractNumber,
		tx.ConsumerID,
		tx.Tenor,
		tx.OTR,
		tx.AdminFee,
		tx.InstallmentAmount,
		tx.InterestAmount,
		tx.AssetName,
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	tx.ID = id
	return nil
}

func (r *transactionRepository) CreateWithTx(sqlTx *sql.Tx, tx *domain.Transaction) error {
	query := `
		INSERT INTO transactions (contract_number, consumer_id, tenor, otr, admin_fee,
		                          installment_amount, interest_amount, asset_name)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := sqlTx.Exec(query,
		tx.ContractNumber,
		tx.ConsumerID,
		tx.Tenor,
		tx.OTR,
		tx.AdminFee,
		tx.InstallmentAmount,
		tx.InterestAmount,
		tx.AssetName,
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	tx.ID = id
	return nil
}

func (r *transactionRepository) FindByContractNumber(contractNumber string) (*domain.Transaction, error) {
	query := `
		SELECT id, contract_number, consumer_id, tenor, otr, admin_fee,
		       installment_amount, interest_amount, asset_name, created_at
		FROM transactions
		WHERE contract_number = ?
	`
	tx := &domain.Transaction{}
	err := r.db.QueryRow(query, contractNumber).Scan(
		&tx.ID,
		&tx.ContractNumber,
		&tx.ConsumerID,
		&tx.Tenor,
		&tx.OTR,
		&tx.AdminFee,
		&tx.InstallmentAmount,
		&tx.InterestAmount,
		&tx.AssetName,
		&tx.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrTransactionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}

	return tx, nil
}

func (r *transactionRepository) FindByConsumerID(consumerID int64) ([]*domain.Transaction, error) {
	query := `
		SELECT id, contract_number, consumer_id, tenor, otr, admin_fee,
		       installment_amount, interest_amount, asset_name, created_at
		FROM transactions
		WHERE consumer_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, consumerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		tx := &domain.Transaction{}
		err := rows.Scan(
			&tx.ID,
			&tx.ContractNumber,
			&tx.ConsumerID,
			&tx.Tenor,
			&tx.OTR,
			&tx.AdminFee,
			&tx.InstallmentAmount,
			&tx.InterestAmount,
			&tx.AssetName,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func (r *transactionRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

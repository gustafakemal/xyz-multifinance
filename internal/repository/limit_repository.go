package repository

import (
	"database/sql"
	"fmt"

	"github.com/username/xyz-multifinance/internal/domain"
)

// LimitRepository defines methods for consumer limit data access
type LimitRepository interface {
	Create(limit *domain.ConsumerLimit) error
	FindByConsumerIDAndTenor(consumerID int64, tenor int) (*domain.ConsumerLimit, error)
	UpdateUsedAmountWithTx(tx *sql.Tx, limitID int64, usedAmount float64, version int) error
	FindAllByConsumerID(consumerID int64) ([]*domain.ConsumerLimit, error)
}

type limitRepository struct {
	db *sql.DB
}

// NewLimitRepository creates a new limit repository
func NewLimitRepository(db *sql.DB) LimitRepository {
	return &limitRepository{db: db}
}

func (r *limitRepository) Create(limit *domain.ConsumerLimit) error {
	query := `
		INSERT INTO consumer_limits (consumer_id, tenor, limit_amount, used_amount, version)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		limit.ConsumerID,
		limit.Tenor,
		limit.LimitAmount,
		limit.UsedAmount,
		limit.Version,
	)
	if err != nil {
		return fmt.Errorf("failed to create consumer limit: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	limit.ID = id
	return nil
}

func (r *limitRepository) FindByConsumerIDAndTenor(consumerID int64, tenor int) (*domain.ConsumerLimit, error) {
	query := `
		SELECT id, consumer_id, tenor, limit_amount, used_amount, version, updated_at
		FROM consumer_limits
		WHERE consumer_id = ? AND tenor = ?
	`
	limit := &domain.ConsumerLimit{}
	err := r.db.QueryRow(query, consumerID, tenor).Scan(
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
		return nil, fmt.Errorf("failed to find consumer limit: %w", err)
	}

	return limit, nil
}

func (r *limitRepository) UpdateUsedAmountWithTx(tx *sql.Tx, limitID int64, usedAmount float64, version int) error {
	// Optimistic locking: update only if version matches
	query := `
		UPDATE consumer_limits
		SET used_amount = ?, version = version + 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND version = ?
	`
	result, err := tx.Exec(query, usedAmount, limitID, version)
	if err != nil {
		return fmt.Errorf("failed to update used amount: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("concurrent update detected: version mismatch")
	}

	return nil
}

func (r *limitRepository) FindAllByConsumerID(consumerID int64) ([]*domain.ConsumerLimit, error) {
	query := `
		SELECT id, consumer_id, tenor, limit_amount, used_amount, version, updated_at
		FROM consumer_limits
		WHERE consumer_id = ?
	`
	rows, err := r.db.Query(query, consumerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query consumer limits: %w", err)
	}
	defer rows.Close()

	var limits []*domain.ConsumerLimit
	for rows.Next() {
		limit := &domain.ConsumerLimit{}
		err := rows.Scan(
			&limit.ID,
			&limit.ConsumerID,
			&limit.Tenor,
			&limit.LimitAmount,
			&limit.UsedAmount,
			&limit.Version,
			&limit.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan consumer limit: %w", err)
		}
		limits = append(limits, limit)
	}

	return limits, nil
}

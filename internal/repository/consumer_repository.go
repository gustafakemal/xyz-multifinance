package repository

import (
	"database/sql"
	"fmt"

	"github.com/username/xyz-multifinance/internal/domain"
)

// ConsumerRepository defines methods for consumer data access
type ConsumerRepository interface {
	Create(consumer *domain.Consumer) error
	FindByID(id int64) (*domain.Consumer, error)
	FindByNIK(nik string) (*domain.Consumer, error)
	Update(consumer *domain.Consumer) error
}

type consumerRepository struct {
	db *sql.DB
}

// NewConsumerRepository creates a new consumer repository
func NewConsumerRepository(db *sql.DB) ConsumerRepository {
	return &consumerRepository{db: db}
}

func (r *consumerRepository) Create(consumer *domain.Consumer) error {
	query := `
		INSERT INTO consumers (nik, full_name, legal_name, birth_place, birth_date, 
		                       salary, ktp_photo, selfie_photo)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		consumer.NIK,
		consumer.FullName,
		consumer.LegalName,
		consumer.BirthPlace,
		consumer.BirthDate,
		consumer.Salary,
		consumer.KTPPhoto,
		consumer.SelfiePhoto,
	)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	consumer.ID = id
	return nil
}

func (r *consumerRepository) FindByID(id int64) (*domain.Consumer, error) {
	query := `
		SELECT id, nik, full_name, legal_name, birth_place, birth_date, 
		       salary, ktp_photo, selfie_photo, created_at
		FROM consumers
		WHERE id = ?
	`
	consumer := &domain.Consumer{}
	err := r.db.QueryRow(query, id).Scan(
		&consumer.ID,
		&consumer.NIK,
		&consumer.FullName,
		&consumer.LegalName,
		&consumer.BirthPlace,
		&consumer.BirthDate,
		&consumer.Salary,
		&consumer.KTPPhoto,
		&consumer.SelfiePhoto,
		&consumer.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrConsumerNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find consumer: %w", err)
	}

	return consumer, nil
}

func (r *consumerRepository) FindByNIK(nik string) (*domain.Consumer, error) {
	query := `
		SELECT id, nik, full_name, legal_name, birth_place, birth_date, 
		       salary, ktp_photo, selfie_photo, created_at
		FROM consumers
		WHERE nik = ?
	`
	consumer := &domain.Consumer{}
	err := r.db.QueryRow(query, nik).Scan(
		&consumer.ID,
		&consumer.NIK,
		&consumer.FullName,
		&consumer.LegalName,
		&consumer.BirthPlace,
		&consumer.BirthDate,
		&consumer.Salary,
		&consumer.KTPPhoto,
		&consumer.SelfiePhoto,
		&consumer.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrConsumerNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find consumer by NIK: %w", err)
	}

	return consumer, nil
}

func (r *consumerRepository) Update(consumer *domain.Consumer) error {
	query := `
		UPDATE consumers
		SET full_name = ?, legal_name = ?, birth_place = ?, birth_date = ?,
		    salary = ?, ktp_photo = ?, selfie_photo = ?
		WHERE id = ?
	`
	result, err := r.db.Exec(query,
		consumer.FullName,
		consumer.LegalName,
		consumer.BirthPlace,
		consumer.BirthDate,
		consumer.Salary,
		consumer.KTPPhoto,
		consumer.SelfiePhoto,
		consumer.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update consumer: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return domain.ErrConsumerNotFound
	}

	return nil
}

package store

import (
	"database/sql"
	"fmt"

	"github.com/ezedu/backend/internal/model"
)

// AccountStore handles account persistence.
type AccountStore struct {
	db *sql.DB
}

func NewAccountStore(db *sql.DB) *AccountStore {
	return &AccountStore{db: db}
}

// Create inserts a new account and returns its ID.
func (s *AccountStore) Create(email, passwordHash, parentName string) (int64, error) {
	result, err := s.db.Exec(
		`INSERT INTO accounts (email, password_hash, parent_name) VALUES (?, ?, ?)`,
		email, passwordHash, parentName,
	)
	if err != nil {
		return 0, fmt.Errorf("create account: %w", err)
	}
	return result.LastInsertId()
}

// GetByEmail retrieves an account by email.
func (s *AccountStore) GetByEmail(email string) (*model.Account, error) {
	a := &model.Account{}
	err := s.db.QueryRow(
		`SELECT id, email, password_hash, parent_name, parent_pin, created_at, updated_at
		 FROM accounts WHERE email = ?`, email,
	).Scan(&a.ID, &a.Email, &a.PasswordHash, &a.ParentName, &a.ParentPIN, &a.CreatedAt, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get account by email: %w", err)
	}
	return a, nil
}

// GetByID retrieves an account by ID.
func (s *AccountStore) GetByID(id int64) (*model.Account, error) {
	a := &model.Account{}
	err := s.db.QueryRow(
		`SELECT id, email, password_hash, parent_name, parent_pin, created_at, updated_at
		 FROM accounts WHERE id = ?`, id,
	).Scan(&a.ID, &a.Email, &a.PasswordHash, &a.ParentName, &a.ParentPIN, &a.CreatedAt, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get account by id: %w", err)
	}
	return a, nil
}

// UpdatePIN updates the parent PIN hash.
func (s *AccountStore) UpdatePIN(accountID int64, pinHash string) error {
	_, err := s.db.Exec(
		`UPDATE accounts SET parent_pin = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		pinHash, accountID,
	)
	return err
}

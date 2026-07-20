package store

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ezedu/backend/internal/model"
)

// SessionStore handles session persistence.
type SessionStore struct {
	db *sql.DB
}

func NewSessionStore(db *sql.DB) *SessionStore {
	return &SessionStore{db: db}
}

// Create creates a new session with a 30-day expiry.
func (s *SessionStore) Create(accountID int64, ipAddress string) (*model.Session, error) {
	token, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("generate session token: %w", err)
	}

	now := time.Now().UTC()
	expiresAt := now.Add(30 * 24 * time.Hour)

	_, err = s.db.Exec(
		`INSERT INTO sessions (id, account_id, created_at, expires_at, ip_address) VALUES (?, ?, ?, ?, ?)`,
		token, accountID, now, expiresAt, ipAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	return &model.Session{
		ID:        token,
		AccountID: accountID,
		CreatedAt: now,
		ExpiresAt: expiresAt,
		IPAddress: ipAddress,
	}, nil
}

// GetByToken retrieves a session by its token, only if not expired.
func (s *SessionStore) GetByToken(token string) (*model.Session, error) {
	sess := &model.Session{}
	err := s.db.QueryRow(
		`SELECT id, account_id, created_at, expires_at, ip_address
		 FROM sessions WHERE id = ? AND expires_at > CURRENT_TIMESTAMP`, token,
	).Scan(&sess.ID, &sess.AccountID, &sess.CreatedAt, &sess.ExpiresAt, &sess.IPAddress)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}
	return sess, nil
}

// Delete removes a session (logout).
func (s *SessionStore) Delete(token string) error {
	_, err := s.db.Exec(`DELETE FROM sessions WHERE id = ?`, token)
	return err
}

// CleanExpired removes expired sessions.
func (s *SessionStore) CleanExpired() error {
	_, err := s.db.Exec(`DELETE FROM sessions WHERE expires_at <= CURRENT_TIMESTAMP`)
	return err
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

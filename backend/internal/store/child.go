package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ezedu/backend/internal/model"
)

// ChildStore handles child profile persistence.
type ChildStore struct {
	db *sql.DB
}

func NewChildStore(db *sql.DB) *ChildStore {
	return &ChildStore{db: db}
}

// Create inserts a new child profile.
func (s *ChildStore) Create(accountID int64, name string, birthYear int, avatarID int) (*model.Child, error) {
	currentYear := time.Now().Year()
	ageGroup := model.AgeGroupFromBirthYear(birthYear, currentYear)

	result, err := s.db.Exec(
		`INSERT INTO children (account_id, name, birth_year, age_group, avatar_id)
		 VALUES (?, ?, ?, ?, ?)`,
		accountID, name, birthYear, ageGroup, avatarID,
	)
	if err != nil {
		return nil, fmt.Errorf("create child: %w", err)
	}

	id, _ := result.LastInsertId()
	return &model.Child{
		ID:           id,
		AccountID:    accountID,
		Name:         name,
		BirthYear:    birthYear,
		AgeGroup:     ageGroup,
		AvatarID:     avatarID,
		XPTotal:      0,
		CurrentLevel: 1,
		StreakDays:    0,
	}, nil
}

// ListByAccount retrieves all children for an account.
func (s *ChildStore) ListByAccount(accountID int64) ([]model.Child, error) {
	rows, err := s.db.Query(
		`SELECT id, account_id, name, birth_year, age_group, avatar_id,
		        xp_total, current_level, streak_days, last_active, daily_limit_min
		 FROM children WHERE account_id = ? ORDER BY id`, accountID,
	)
	if err != nil {
		return nil, fmt.Errorf("list children: %w", err)
	}
	defer rows.Close()

	var children []model.Child
	for rows.Next() {
		var c model.Child
		if err := rows.Scan(
			&c.ID, &c.AccountID, &c.Name, &c.BirthYear, &c.AgeGroup,
			&c.AvatarID, &c.XPTotal, &c.CurrentLevel, &c.StreakDays,
			&c.LastActive, &c.DailyLimitMin,
		); err != nil {
			return nil, fmt.Errorf("scan child: %w", err)
		}
		children = append(children, c)
	}
	return children, nil
}

// GetByID retrieves a child by ID, scoped to an account.
func (s *ChildStore) GetByID(id, accountID int64) (*model.Child, error) {
	c := &model.Child{}
	err := s.db.QueryRow(
		`SELECT id, account_id, name, birth_year, age_group, avatar_id,
		        xp_total, current_level, streak_days, last_active, daily_limit_min
		 FROM children WHERE id = ? AND account_id = ?`, id, accountID,
	).Scan(
		&c.ID, &c.AccountID, &c.Name, &c.BirthYear, &c.AgeGroup,
		&c.AvatarID, &c.XPTotal, &c.CurrentLevel, &c.StreakDays,
		&c.LastActive, &c.DailyLimitMin,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get child: %w", err)
	}
	return c, nil
}

// Update modifies a child's name, birth year, and avatar.
func (s *ChildStore) Update(id, accountID int64, name string, birthYear, avatarID int) error {
	currentYear := time.Now().Year()
	ageGroup := model.AgeGroupFromBirthYear(birthYear, currentYear)

	_, err := s.db.Exec(
		`UPDATE children SET name = ?, birth_year = ?, age_group = ?, avatar_id = ?
		 WHERE id = ? AND account_id = ?`,
		name, birthYear, ageGroup, avatarID, id, accountID,
	)
	return err
}

// Delete removes a child profile.
func (s *ChildStore) Delete(id, accountID int64) error {
	_, err := s.db.Exec(`DELETE FROM children WHERE id = ? AND account_id = ?`, id, accountID)
	return err
}

// CountByAccount returns the number of child profiles for an account.
func (s *ChildStore) CountByAccount(accountID int64) (int, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM children WHERE account_id = ?`, accountID).Scan(&count)
	return count, err
}

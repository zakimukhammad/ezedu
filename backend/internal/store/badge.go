package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ezedu/backend/internal/model"
)

// BadgeStore handles badge queries and achievement evaluation.
type BadgeStore struct {
	db *sql.DB
}

func NewBadgeStore(db *sql.DB) *BadgeStore {
	return &BadgeStore{db: db}
}

// EarnedBadgeResponse represents a badge with child's earned status.
type EarnedBadgeResponse struct {
	model.Badge
	Earned   bool       `json:"earned"`
	EarnedAt *time.Time `json:"earned_at,omitempty"`
}

// ListChildBadges returns all badges in system with status for a specific child.
func (s *BadgeStore) ListChildBadges(childID int64) ([]EarnedBadgeResponse, error) {
	rows, err := s.db.Query(
		`SELECT b.id, b.slug, b.name, b.description, b.icon, b.category_id, b.criteria_json,
		        cb.earned_at
		 FROM badges b
		 LEFT JOIN child_badges cb ON b.id = cb.badge_id AND cb.child_id = ?
		 ORDER BY b.id ASC`,
		childID,
	)
	if err != nil {
		return nil, fmt.Errorf("list child badges: %w", err)
	}
	defer rows.Close()

	var list []EarnedBadgeResponse
	for rows.Next() {
		var eb EarnedBadgeResponse
		var earnedAt sql.NullTime
		if err := rows.Scan(
			&eb.ID, &eb.Slug, &eb.Name, &eb.Description, &eb.Icon,
			&eb.CategoryID, &eb.CriteriaJSON, &earnedAt,
		); err != nil {
			return nil, fmt.Errorf("scan badge: %w", err)
		}
		if earnedAt.Valid {
			eb.Earned = true
			eb.EarnedAt = &earnedAt.Time
		}
		list = append(list, eb)
	}

	return list, nil
}

// EvaluateAndAwardBadges checks if a child has fulfilled criteria for new badges and awards them.
func (s *BadgeStore) EvaluateAndAwardBadges(childID int64) ([]model.Badge, error) {
	// Get unearned badges
	allBadges, err := s.ListChildBadges(childID)
	if err != nil {
		return nil, err
	}

	// Fetch child stats
	var totalCompleted int
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM child_progress WHERE child_id = ? AND status = 'completed'`, childID).Scan(&totalCompleted)

	var mathCompleted int
	_ = s.db.QueryRow(
		`SELECT COUNT(*) FROM child_progress cp 
		 JOIN lessons l ON cp.lesson_id = l.id 
		 JOIN categories c ON l.category_id = c.id 
		 WHERE cp.child_id = ? AND cp.status = 'completed' AND c.slug = 'math'`,
		childID,
	).Scan(&mathCompleted)

	var codingCompleted int
	_ = s.db.QueryRow(
		`SELECT COUNT(*) FROM child_progress cp 
		 JOIN lessons l ON cp.lesson_id = l.id 
		 JOIN categories c ON l.category_id = c.id 
		 WHERE cp.child_id = ? AND cp.status = 'completed' AND c.slug = 'coding'`,
		childID,
	).Scan(&codingCompleted)

	var perfectScores int
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM child_progress WHERE child_id = ? AND score >= max_possible AND max_possible > 0`, childID).Scan(&perfectScores)

	var xpTotal, currentLevel, streakDays int
	_ = s.db.QueryRow(`SELECT xp_total, current_level, streak_days FROM children WHERE id = ?`, childID).Scan(&xpTotal, &currentLevel, &streakDays)

	// Fetch category completion counts helper
	getCategoryCompleted := func(slug string) int {
		var count int
		_ = s.db.QueryRow(
			`SELECT COUNT(*) FROM child_progress cp 
			 JOIN lessons l ON cp.lesson_id = l.id 
			 JOIN categories c ON l.category_id = c.id 
			 WHERE cp.child_id = ? AND cp.status = 'completed' AND c.slug = ?`,
			childID, slug,
		).Scan(&count)
		return count
	}

	var newlyAwarded []model.Badge

	for _, b := range allBadges {
		if b.Earned {
			continue
		}

		var crit struct {
			Type     string `json:"type"`
			Value    int    `json:"value"`
			Category string `json:"category"`
		}
		if err := json.Unmarshal([]byte(b.CriteriaJSON), &crit); err != nil {
			continue
		}

		shouldAward := false
		switch crit.Type {
		case "lessons_completed":
			if totalCompleted >= crit.Value {
				shouldAward = true
			}
		case "category_lessons":
			if getCategoryCompleted(crit.Category) >= crit.Value {
				shouldAward = true
			}
		case "perfect_score":
			if perfectScores >= crit.Value {
				shouldAward = true
			}
		case "streak":
			if streakDays >= crit.Value {
				shouldAward = true
			}
		case "xp_total":
			if xpTotal >= crit.Value {
				shouldAward = true
			}
		case "level":
			if currentLevel >= crit.Value {
				shouldAward = true
			}
		}

		if shouldAward {
			_, err := s.db.Exec(
				`INSERT OR IGNORE INTO child_badges (child_id, badge_id, earned_at) VALUES (?, ?, CURRENT_TIMESTAMP)`,
				childID, b.ID,
			)
			if err == nil {
				newlyAwarded = append(newlyAwarded, b.Badge)
			}
		}
	}

	return newlyAwarded, nil
}

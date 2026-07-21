package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ezedu/backend/internal/model"
)

// ProgressStore handles child lesson progress and activity results.
type ProgressStore struct {
	db *sql.DB
}

func NewProgressStore(db *sql.DB) *ProgressStore {
	return &ProgressStore{db: db}
}

// GetChildLessonProgress retrieves progress for a specific child and lesson.
func (s *ProgressStore) GetChildLessonProgress(childID, lessonID int64) (*model.ChildProgress, error) {
	p := &model.ChildProgress{}
	err := s.db.QueryRow(
		`SELECT id, child_id, lesson_id, status, score, max_possible, attempts, started_at, completed_at, time_spent_sec
		 FROM child_progress WHERE child_id = ? AND lesson_id = ?`,
		childID, lessonID,
	).Scan(
		&p.ID, &p.ChildID, &p.LessonID, &p.Status, &p.Score, &p.MaxPossible,
		&p.Attempts, &p.StartedAt, &p.CompletedAt, &p.TimeSpentSec,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get child progress: %w", err)
	}
	return p, nil
}

// ListChildProgressByChild retrieves all progress entries for a child.
func (s *ProgressStore) ListChildProgressByChild(childID int64) ([]model.ChildProgress, error) {
	rows, err := s.db.Query(
		`SELECT id, child_id, lesson_id, status, score, max_possible, attempts, started_at, completed_at, time_spent_sec
		 FROM child_progress WHERE child_id = ?`,
		childID,
	)
	if err != nil {
		return nil, fmt.Errorf("list progress by child: %w", err)
	}
	defer rows.Close()

	var list []model.ChildProgress
	for rows.Next() {
		var p model.ChildProgress
		if err := rows.Scan(
			&p.ID, &p.ChildID, &p.LessonID, &p.Status, &p.Score, &p.MaxPossible,
			&p.Attempts, &p.StartedAt, &p.CompletedAt, &p.TimeSpentSec,
		); err != nil {
			return nil, fmt.Errorf("scan child progress: %w", err)
		}
		list = append(list, p)
	}
	return list, nil
}

// RecordActivityResult saves an activity submission result.
func (s *ProgressStore) RecordActivityResult(childID, activityID int64, answerJSON string, isCorrect bool, score, attemptNumber int) error {
	_, err := s.db.Exec(
		`INSERT INTO activity_results (child_id, activity_id, answer_json, is_correct, score, attempt_number)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		childID, activityID, answerJSON, isCorrect, score, attemptNumber,
	)
	if err != nil {
		return fmt.Errorf("record activity result: %w", err)
	}
	return nil
}

// CompleteLesson updates or creates the progress row for a lesson and awards XP to the child.
func (s *ProgressStore) CompleteLesson(childID, lessonID int64, finalScore, maxScore, timeSpentSec, xpReward int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// Upsert lesson progress
	_, err = tx.Exec(
		`INSERT INTO child_progress (child_id, lesson_id, status, score, max_possible, attempts, started_at, completed_at, time_spent_sec)
		 VALUES (?, ?, 'completed', ?, ?, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, ?)
		 ON CONFLICT(child_id, lesson_id) DO UPDATE SET
		   status = 'completed',
		   score = MAX(score, excluded.score),
		   max_possible = excluded.max_possible,
		   attempts = attempts + 1,
		   completed_at = CURRENT_TIMESTAMP,
		   time_spent_sec = time_spent_sec + excluded.time_spent_sec`,
		childID, lessonID, finalScore, maxScore, timeSpentSec,
	)
	if err != nil {
		return fmt.Errorf("upsert lesson progress: %w", err)
	}

	// Award XP to child profile
	_, err = tx.Exec(
		`UPDATE children 
		 SET xp_total = xp_total + ?, 
		     current_level = 1 + ((xp_total + ?) / 100),
		     last_active = CURRENT_DATE
		 WHERE id = ?`,
		xpReward, xpReward, childID,
	)
	if err != nil {
		return fmt.Errorf("update child xp: %w", err)
	}

	return tx.Commit()
}

// UpdateStreak recalculates and updates a child's streak_days based on last_active.
func (s *ProgressStore) UpdateStreak(childID int64) error {
	var lastActive *string
	_ = s.db.QueryRow(`SELECT last_active FROM children WHERE id = ?`, childID).Scan(&lastActive)

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	if lastActive == nil || *lastActive == "" {
		// First activity ever — set streak to 1
		_, err := s.db.Exec(`UPDATE children SET streak_days = 1, last_active = ? WHERE id = ?`, today, childID)
		return err
	}

	switch *lastActive {
	case today:
		// Already active today — no change
		return nil
	case yesterday:
		// Consecutive day — increment streak
		_, err := s.db.Exec(`UPDATE children SET streak_days = streak_days + 1, last_active = ? WHERE id = ?`, today, childID)
		return err
	default:
		// Streak broken — reset to 1
		_, err := s.db.Exec(`UPDATE children SET streak_days = 1, last_active = ? WHERE id = ?`, today, childID)
		return err
	}
}

// ChildProgressSummary holds aggregated progress data for a child.
type ChildProgressSummary struct {
	TotalLessonsCompleted int                     `json:"total_lessons_completed"`
	TotalScore            int                     `json:"total_score"`
	TotalMaxPossible      int                     `json:"total_max_possible"`
	TotalTimeSpentSec     int                     `json:"total_time_spent_sec"`
	CategoryProgress      []CategoryProgressEntry `json:"category_progress"`
	RecentActivity        []RecentActivityEntry   `json:"recent_activity"`
}

type CategoryProgressEntry struct {
	CategorySlug   string `json:"category_slug"`
	CategoryName   string `json:"category_name"`
	CategoryColor  string `json:"category_color"`
	Completed      int    `json:"completed"`
	TotalAvailable int    `json:"total_available"`
	Score          int    `json:"score"`
	MaxPossible    int    `json:"max_possible"`
}

type RecentActivityEntry struct {
	LessonID    int64  `json:"lesson_id"`
	LessonTitle string `json:"lesson_title"`
	Status      string `json:"status"`
	Score       int    `json:"score"`
	MaxPossible int    `json:"max_possible"`
	CompletedAt string `json:"completed_at"`
}

// GetChildProgressSummary returns an aggregated progress summary for a child.
func (s *ProgressStore) GetChildProgressSummary(childID int64, ageGroup string) (*ChildProgressSummary, error) {
	summary := &ChildProgressSummary{}

	// Aggregate totals
	_ = s.db.QueryRow(
		`SELECT COALESCE(COUNT(*), 0), COALESCE(SUM(score), 0), COALESCE(SUM(max_possible), 0), COALESCE(SUM(time_spent_sec), 0)
		 FROM child_progress WHERE child_id = ? AND status = 'completed'`,
		childID,
	).Scan(&summary.TotalLessonsCompleted, &summary.TotalScore, &summary.TotalMaxPossible, &summary.TotalTimeSpentSec)

	// Category-wise breakdown
	catRows, err := s.db.Query(
		`SELECT c.slug, c.name, c.color,
		        COALESCE(SUM(CASE WHEN cp.status = 'completed' THEN 1 ELSE 0 END), 0) as completed,
		        COUNT(l.id) as total_available,
		        COALESCE(SUM(CASE WHEN cp.status = 'completed' THEN cp.score ELSE 0 END), 0) as score,
		        COALESCE(SUM(CASE WHEN cp.status = 'completed' THEN cp.max_possible ELSE 0 END), 0) as max_possible
		 FROM categories c
		 JOIN lessons l ON l.category_id = c.id AND l.age_group = ?
		 LEFT JOIN child_progress cp ON cp.lesson_id = l.id AND cp.child_id = ?
		 GROUP BY c.id
		 ORDER BY c.sort_order`,
		ageGroup, childID,
	)
	if err != nil {
		return nil, fmt.Errorf("category progress: %w", err)
	}
	defer catRows.Close()

	for catRows.Next() {
		var entry CategoryProgressEntry
		if err := catRows.Scan(
			&entry.CategorySlug, &entry.CategoryName, &entry.CategoryColor,
			&entry.Completed, &entry.TotalAvailable, &entry.Score, &entry.MaxPossible,
		); err != nil {
			return nil, fmt.Errorf("scan category progress: %w", err)
		}
		summary.CategoryProgress = append(summary.CategoryProgress, entry)
	}

	// Recent activity (last 10)
	recentRows, err := s.db.Query(
		`SELECT cp.lesson_id, l.title, cp.status, cp.score, cp.max_possible, COALESCE(cp.completed_at, cp.started_at, '') as completed_at
		 FROM child_progress cp
		 JOIN lessons l ON cp.lesson_id = l.id
		 WHERE cp.child_id = ?
		 ORDER BY COALESCE(cp.completed_at, cp.started_at) DESC
		 LIMIT 10`,
		childID,
	)
	if err != nil {
		return nil, fmt.Errorf("recent activity: %w", err)
	}
	defer recentRows.Close()

	for recentRows.Next() {
		var entry RecentActivityEntry
		if err := recentRows.Scan(
			&entry.LessonID, &entry.LessonTitle, &entry.Status, &entry.Score, &entry.MaxPossible, &entry.CompletedAt,
		); err != nil {
			return nil, fmt.Errorf("scan recent activity: %w", err)
		}
		summary.RecentActivity = append(summary.RecentActivity, entry)
	}

	return summary, nil
}

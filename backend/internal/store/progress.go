package store

import (
	"database/sql"
	"fmt"

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

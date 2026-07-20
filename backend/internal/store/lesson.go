package store

import (
	"database/sql"
	"fmt"

	"github.com/ezedu/backend/internal/model"
)

// LessonStore handles lesson and activity queries.
type LessonStore struct {
	db *sql.DB
}

func NewLessonStore(db *sql.DB) *LessonStore {
	return &LessonStore{db: db}
}

// ListByCategoryAndAgeGroup retrieves lessons for a specific category and age group.
func (s *LessonStore) ListByCategoryAndAgeGroup(categoryID int64, ageGroup string) ([]model.Lesson, error) {
	rows, err := s.db.Query(
		`SELECT id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward
		 FROM lessons 
		 WHERE category_id = ? AND age_group = ?
		 ORDER BY level ASC, sort_order ASC`,
		categoryID, ageGroup,
	)
	if err != nil {
		return nil, fmt.Errorf("list lessons: %w", err)
	}
	defer rows.Close()

	var lessons []model.Lesson
	for rows.Next() {
		var l model.Lesson
		if err := rows.Scan(
			&l.ID, &l.CategoryID, &l.AgeGroup, &l.Level, &l.SortOrder,
			&l.Title, &l.Description, &l.ContentJSON, &l.EstimatedMinutes, &l.XPReward,
		); err != nil {
			return nil, fmt.Errorf("scan lesson: %w", err)
		}
		lessons = append(lessons, l)
	}

	return lessons, nil
}

// GetByID retrieves a lesson by its ID.
func (s *LessonStore) GetByID(id int64) (*model.Lesson, error) {
	l := &model.Lesson{}
	err := s.db.QueryRow(
		`SELECT id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward
		 FROM lessons WHERE id = ?`, id,
	).Scan(
		&l.ID, &l.CategoryID, &l.AgeGroup, &l.Level, &l.SortOrder,
		&l.Title, &l.Description, &l.ContentJSON, &l.EstimatedMinutes, &l.XPReward,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get lesson by id: %w", err)
	}
	return l, nil
}

// ListActivitiesByLessonID retrieves all activities for a lesson.
func (s *LessonStore) ListActivitiesByLessonID(lessonID int64) ([]model.Activity, error) {
	rows, err := s.db.Query(
		`SELECT id, lesson_id, type, sort_order, question_json, max_score
		 FROM activities 
		 WHERE lesson_id = ?
		 ORDER BY sort_order ASC`,
		lessonID,
	)
	if err != nil {
		return nil, fmt.Errorf("list activities: %w", err)
	}
	defer rows.Close()

	var activities []model.Activity
	for rows.Next() {
		var a model.Activity
		if err := rows.Scan(
			&a.ID, &a.LessonID, &a.Type, &a.SortOrder, &a.QuestionJSON, &a.MaxScore,
		); err != nil {
			return nil, fmt.Errorf("scan activity: %w", err)
		}
		activities = append(activities, a)
	}

	return activities, nil
}

// GetActivityByID retrieves a single activity.
func (s *LessonStore) GetActivityByID(id int64) (*model.Activity, error) {
	a := &model.Activity{}
	err := s.db.QueryRow(
		`SELECT id, lesson_id, type, sort_order, question_json, max_score
		 FROM activities WHERE id = ?`, id,
	).Scan(&a.ID, &a.LessonID, &a.Type, &a.SortOrder, &a.QuestionJSON, &a.MaxScore)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get activity by id: %w", err)
	}
	return a, nil
}

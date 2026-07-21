package model

import "time"

// Account represents a parent/guardian account.
type Account struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	ParentName   string    `json:"parent_name"`
	ParentPIN    string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Child represents a child profile within an account.
type Child struct {
	ID            int64   `json:"id"`
	AccountID     int64   `json:"account_id"`
	Name          string  `json:"name"`
	BirthYear     int     `json:"birth_year"`
	AgeGroup      string  `json:"age_group"`
	AvatarID      int     `json:"avatar_id"`
	XPTotal       int     `json:"xp_total"`
	CurrentLevel  int     `json:"current_level"`
	StreakDays    int     `json:"streak_days"`
	LastActive    *string `json:"last_active,omitempty"`
	DailyLimitMin *int    `json:"daily_limit_min,omitempty"`
}

// Category represents a learning category (e.g., Matematika, Sains).
type Category struct {
	ID          int64  `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	SortOrder   int    `json:"sort_order"`
}

// Lesson represents a single lesson within a category.
type Lesson struct {
	ID               int64  `json:"id"`
	CategoryID       int64  `json:"category_id"`
	CategorySlug     string `json:"category_slug,omitempty"`
	AgeGroup         string `json:"age_group"`
	Level            int    `json:"level"`
	SortOrder        int    `json:"sort_order"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	ContentJSON      string `json:"content_json"`
	EstimatedMinutes int    `json:"estimated_minutes"`
	XPReward         int    `json:"xp_reward"`
}

// Activity represents an interactive exercise within a lesson.
type Activity struct {
	ID           int64  `json:"id"`
	LessonID     int64  `json:"lesson_id"`
	Type         string `json:"type"`
	SortOrder    int    `json:"sort_order"`
	QuestionJSON string `json:"question_json"`
	MaxScore     int    `json:"max_score"`
}

// ChildProgress tracks a child's progress on a lesson.
type ChildProgress struct {
	ID           int64   `json:"id"`
	ChildID      int64   `json:"child_id"`
	LessonID     int64   `json:"lesson_id"`
	Status       string  `json:"status"`
	Score        int     `json:"score"`
	MaxPossible  int     `json:"max_possible"`
	Attempts     int     `json:"attempts"`
	StartedAt    *string `json:"started_at,omitempty"`
	CompletedAt  *string `json:"completed_at,omitempty"`
	TimeSpentSec int     `json:"time_spent_sec"`
}

// ActivityResult records the result of a child answering an activity.
type ActivityResult struct {
	ID            int64     `json:"id"`
	ChildID       int64     `json:"child_id"`
	ActivityID    int64     `json:"activity_id"`
	AnswerJSON    string    `json:"answer_json"`
	IsCorrect     bool      `json:"is_correct"`
	Score         int       `json:"score"`
	AttemptNumber int       `json:"attempt_number"`
	AnsweredAt    time.Time `json:"answered_at"`
}

// Badge represents an earnable achievement.
type Badge struct {
	ID           int64  `json:"id"`
	Slug         string `json:"slug"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Icon         string `json:"icon"`
	CategoryID   *int64 `json:"category_id,omitempty"`
	CriteriaJSON string `json:"criteria_json"`
}

// Session represents an active login session.
type Session struct {
	ID        string    `json:"id"`
	AccountID int64     `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	IPAddress string    `json:"ip_address"`
}

// AgeGroupFromBirthYear derives the age group based on birth year and current year.
func AgeGroupFromBirthYear(birthYear, currentYear int) string {
	age := currentYear - birthYear
	switch {
	case age >= 4 && age <= 6:
		return "explorers"
	case age >= 7 && age <= 9:
		return "builders"
	case age >= 10 && age <= 12:
		return "challengers"
	case age < 4:
		return "toddlers"
	default:
		return "challengers"
	}
}

// DailyChallenge represents a daily challenge for an age group.
type DailyChallenge struct {
	ID            int64  `json:"id"`
	ChallengeDate string `json:"challenge_date"`
	AgeGroup      string `json:"age_group"`
	ActivityType  string `json:"activity_type"`
	QuestionJSON  string `json:"question_json"`
	MaxScore      int    `json:"max_score"`
}

// DailyChallengeResult records a child's result for a daily challenge.
type DailyChallengeResult struct {
	ID          int64  `json:"id"`
	ChildID     int64  `json:"child_id"`
	ChallengeID int64  `json:"challenge_id"`
	Score       int    `json:"score"`
	CompletedAt string `json:"completed_at"`
}

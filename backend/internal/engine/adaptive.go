package engine

import (
	"database/sql"
	"fmt"

	"github.com/ezedu/backend/internal/model"
)

const (
	LowScoreThreshold   = 50.0 // < 50% score triggers low counter
	HighScoreThreshold  = 90.0 // > 90% score triggers high counter
	ConsecutiveRequired = 3    // 3 consecutive required to trigger recommendation
	MinLevel            = 1
	MaxLevel            = 5
)

// AdaptiveEngine handles difficulty evaluation and adjustment state.
type AdaptiveEngine struct {
	db *sql.DB
}

func NewAdaptiveEngine(db *sql.DB) *AdaptiveEngine {
	return &AdaptiveEngine{db: db}
}

// EvaluateAfterActivity checks a child's activity score and updates difficulty counters.
func (e *AdaptiveEngine) EvaluateAfterActivity(childID, activityID int64, score, maxScore int) (*model.DifficultyAdjustment, error) {
	if childID <= 0 || activityID <= 0 || maxScore <= 0 {
		return nil, nil
	}

	// 1. Check if child is a toddler (toddlers are excluded from adaptive difficulty)
	var ageGroup string
	err := e.db.QueryRow(`SELECT age_group FROM children WHERE id = ?`, childID).Scan(&ageGroup)
	if err != nil {
		return nil, fmt.Errorf("evaluate child age_group: %w", err)
	}
	if ageGroup == "toddlers" {
		return nil, nil
	}

	// 2. Get activity's lesson details (category_id, level, category_slug, category_name)
	var categoryID int64
	var lessonLevel int
	var catSlug, catName string
	err = e.db.QueryRow(
		`SELECT l.category_id, l.level, c.slug, c.name
		 FROM activities a
		 JOIN lessons l ON a.lesson_id = l.id
		 JOIN categories c ON l.category_id = c.id
		 WHERE a.id = ?`, activityID,
	).Scan(&categoryID, &lessonLevel, &catSlug, &catName)
	if err != nil {
		return nil, fmt.Errorf("evaluate activity details: %w", err)
	}

	// 3. Get or initialize difficulty_adjustments record
	adj := &model.DifficultyAdjustment{
		ChildID:      childID,
		CategoryID:   categoryID,
		CategorySlug: catSlug,
		CategoryName: catName,
		CurrentLevel: lessonLevel,
	}

	var recLevel sql.NullInt64
	var recType string
	err = e.db.QueryRow(
		`SELECT id, current_level, recommended_level, recommendation, consecutive_low, consecutive_high
		 FROM difficulty_adjustments
		 WHERE child_id = ? AND category_id = ?`,
		childID, categoryID,
	).Scan(
		&adj.ID, &adj.CurrentLevel, &recLevel, &recType,
		&adj.ConsecutiveLow, &adj.ConsecutiveHigh,
	)

	if err == sql.ErrNoRows {
		// Initialize new record
		res, err := e.db.Exec(
			`INSERT INTO difficulty_adjustments (child_id, category_id, current_level)
			 VALUES (?, ?, ?)`,
			childID, categoryID, lessonLevel,
		)
		if err != nil {
			return nil, fmt.Errorf("init difficulty adjustment: %w", err)
		}
		adj.ID, _ = res.LastInsertId()
		adj.CurrentLevel = lessonLevel
	} else if err != nil {
		return nil, fmt.Errorf("get difficulty adjustment: %w", err)
	} else {
		if recLevel.Valid {
			v := int(recLevel.Int64)
			adj.RecommendedLevel = &v
		}
		adj.Recommendation = recType
	}

	// 4. Calculate score percentage
	scorePct := (float64(score) / float64(maxScore)) * 100.0

	// 5. Update consecutive counters
	if scorePct < LowScoreThreshold {
		adj.ConsecutiveLow++
		adj.ConsecutiveHigh = 0
	} else if scorePct > HighScoreThreshold {
		adj.ConsecutiveHigh++
		adj.ConsecutiveLow = 0
	} else {
		adj.ConsecutiveLow = 0
		adj.ConsecutiveHigh = 0
	}

	// 6. Check for recommendation triggers
	var nextRecLevel *int = nil
	nextRecType := ""

	if adj.ConsecutiveLow >= ConsecutiveRequired && adj.CurrentLevel > MinLevel {
		nextRecType = "easier"
		target := adj.CurrentLevel - 1
		nextRecLevel = &target
	} else if adj.ConsecutiveHigh >= ConsecutiveRequired && adj.CurrentLevel < MaxLevel {
		nextRecType = "harder"
		target := adj.CurrentLevel + 1
		nextRecLevel = &target
	}

	adj.Recommendation = nextRecType
	adj.RecommendedLevel = nextRecLevel

	// 7. Update database
	var dbRecLevel interface{} = nil
	if nextRecLevel != nil {
		dbRecLevel = *nextRecLevel
	}

	_, err = e.db.Exec(
		`UPDATE difficulty_adjustments
		 SET current_level = ?, recommended_level = ?, recommendation = ?,
		     consecutive_low = ?, consecutive_high = ?, last_evaluated_at = CURRENT_TIMESTAMP
		 WHERE id = ?`,
		adj.CurrentLevel, dbRecLevel, adj.Recommendation,
		adj.ConsecutiveLow, adj.ConsecutiveHigh, adj.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("update difficulty adjustment: %w", err)
	}

	return adj, nil
}

// AcceptRecommendation applies the recommended level change.
func (e *AdaptiveEngine) AcceptRecommendation(childID, categoryID int64) (*model.DifficultyAdjustment, error) {
	var id int64
	var currentLvl int
	var recLevel sql.NullInt64
	var recType string

	err := e.db.QueryRow(
		`SELECT id, current_level, recommended_level, recommendation
		 FROM difficulty_adjustments
		 WHERE child_id = ? AND category_id = ?`,
		childID, categoryID,
	).Scan(&id, &currentLvl, &recLevel, &recType)

	if err != nil {
		return nil, fmt.Errorf("accept recommendation: %w", err)
	}

	if !recLevel.Valid || recType == "" {
		return nil, fmt.Errorf("no pending recommendation to accept")
	}

	newLevel := int(recLevel.Int64)

	_, err = e.db.Exec(
		`UPDATE difficulty_adjustments
		 SET current_level = ?, recommended_level = NULL, recommendation = '',
		     consecutive_low = 0, consecutive_high = 0, last_evaluated_at = CURRENT_TIMESTAMP
		 WHERE id = ?`,
		newLevel, id,
	)
	if err != nil {
		return nil, fmt.Errorf("apply recommendation: %w", err)
	}

	return e.GetCategoryDifficulty(childID, categoryID)
}

// DismissRecommendation clears the pending recommendation.
func (e *AdaptiveEngine) DismissRecommendation(childID, categoryID int64) error {
	_, err := e.db.Exec(
		`UPDATE difficulty_adjustments
		 SET recommended_level = NULL, recommendation = '',
		     consecutive_low = 0, consecutive_high = 0, last_evaluated_at = CURRENT_TIMESTAMP
		 WHERE child_id = ? AND category_id = ?`,
		childID, categoryID,
	)
	return err
}

// GetCategoryDifficulty returns difficulty state for a child in a specific category.
func (e *AdaptiveEngine) GetCategoryDifficulty(childID, categoryID int64) (*model.DifficultyAdjustment, error) {
	adj := &model.DifficultyAdjustment{
		ChildID:    childID,
		CategoryID: categoryID,
	}

	var recLevel sql.NullInt64
	var recType string
	var lastEval sql.NullString

	err := e.db.QueryRow(
		`SELECT d.id, d.current_level, d.recommended_level, d.recommendation,
		        d.consecutive_low, d.consecutive_high, d.last_evaluated_at,
		        c.slug, c.name
		 FROM difficulty_adjustments d
		 JOIN categories c ON d.category_id = c.id
		 WHERE d.child_id = ? AND d.category_id = ?`,
		childID, categoryID,
	).Scan(
		&adj.ID, &adj.CurrentLevel, &recLevel, &recType,
		&adj.ConsecutiveLow, &adj.ConsecutiveHigh, &lastEval,
		&adj.CategorySlug, &adj.CategoryName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get category difficulty: %w", err)
	}

	if recLevel.Valid {
		v := int(recLevel.Int64)
		adj.RecommendedLevel = &v
	}
	adj.Recommendation = recType
	if lastEval.Valid {
		adj.LastEvaluatedAt = &lastEval.String
	}

	return adj, nil
}

// GetDifficultyStates returns difficulty adjustments for all categories for a child.
func (e *AdaptiveEngine) GetDifficultyStates(childID int64) ([]model.DifficultyAdjustment, error) {
	rows, err := e.db.Query(
		`SELECT d.id, d.category_id, d.current_level, d.recommended_level, d.recommendation,
		        d.consecutive_low, d.consecutive_high, d.last_evaluated_at,
		        c.slug, c.name
		 FROM difficulty_adjustments d
		 JOIN categories c ON d.category_id = c.id
		 WHERE d.child_id = ?
		 ORDER BY c.sort_order ASC`,
		childID,
	)
	if err != nil {
		return nil, fmt.Errorf("get difficulty states: %w", err)
	}
	defer rows.Close()

	var list []model.DifficultyAdjustment
	for rows.Next() {
		var adj model.DifficultyAdjustment
		adj.ChildID = childID
		var recLevel sql.NullInt64
		var recType string
		var lastEval sql.NullString

		if err := rows.Scan(
			&adj.ID, &adj.CategoryID, &adj.CurrentLevel, &recLevel, &recType,
			&adj.ConsecutiveLow, &adj.ConsecutiveHigh, &lastEval,
			&adj.CategorySlug, &adj.CategoryName,
		); err != nil {
			return nil, fmt.Errorf("scan difficulty state: %w", err)
		}

		if recLevel.Valid {
			v := int(recLevel.Int64)
			adj.RecommendedLevel = &v
		}
		adj.Recommendation = recType
		if lastEval.Valid {
			adj.LastEvaluatedAt = &lastEval.String
		}

		list = append(list, adj)
	}

	return list, nil
}

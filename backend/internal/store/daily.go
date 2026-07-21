package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/ezedu/backend/internal/model"
)

// DailyStore handles daily challenge queries.
type DailyStore struct {
	db *sql.DB
}

func NewDailyStore(db *sql.DB) *DailyStore {
	return &DailyStore{db: db}
}

// GetTodayChallenge returns today's challenge for the given age group, auto-generating if needed.
func (s *DailyStore) GetTodayChallenge(ageGroup string) (*model.DailyChallenge, error) {
	today := time.Now().Format("2006-01-02")

	c := &model.DailyChallenge{}
	err := s.db.QueryRow(
		`SELECT id, challenge_date, age_group, activity_type, question_json, max_score
		 FROM daily_challenges WHERE challenge_date = ? AND age_group = ?`,
		today, ageGroup,
	).Scan(&c.ID, &c.ChallengeDate, &c.AgeGroup, &c.ActivityType, &c.QuestionJSON, &c.MaxScore)

	if err == sql.ErrNoRows {
		// Auto-generate a new daily challenge
		c, err = s.generateChallenge(today, ageGroup)
		if err != nil {
			return nil, fmt.Errorf("generate daily challenge: %w", err)
		}
		return c, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get today challenge: %w", err)
	}
	return c, nil
}

// SubmitResult records a child's daily challenge result.
func (s *DailyStore) SubmitResult(childID, challengeID int64, score int) error {
	_, err := s.db.Exec(
		`INSERT OR IGNORE INTO daily_challenge_results (child_id, challenge_id, score) VALUES (?, ?, ?)`,
		childID, challengeID, score,
	)
	return err
}

// HasCompleted checks if a child has already completed today's challenge.
func (s *DailyStore) HasCompleted(childID, challengeID int64) bool {
	var id int64
	err := s.db.QueryRow(
		`SELECT id FROM daily_challenge_results WHERE child_id = ? AND challenge_id = ?`,
		childID, challengeID,
	).Scan(&id)
	return err == nil
}

// GetStreak returns the number of consecutive days a child has completed daily challenges.
func (s *DailyStore) GetStreak(childID int64) int {
	rows, err := s.db.Query(
		`SELECT dc.challenge_date FROM daily_challenge_results dcr
		 JOIN daily_challenges dc ON dc.id = dcr.challenge_id
		 WHERE dcr.child_id = ?
		 ORDER BY dc.challenge_date DESC`,
		childID,
	)
	if err != nil {
		return 0
	}
	defer rows.Close()

	streak := 0
	expectedDate := time.Now()
	for rows.Next() {
		var dateStr string
		if err := rows.Scan(&dateStr); err != nil {
			break
		}
		d, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			break
		}

		expectedFormatted := expectedDate.Format("2006-01-02")
		if dateStr == expectedFormatted {
			streak++
			expectedDate = expectedDate.AddDate(0, 0, -1)
		} else if d.Before(expectedDate) {
			break
		}
	}
	return streak
}

// generateChallenge creates a randomized daily challenge for the given date and age group.
func (s *DailyStore) generateChallenge(date, ageGroup string) (*model.DailyChallenge, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	type challengeTemplate struct {
		activityType string
		questionJSON string
		maxScore     int
	}

	// Pool of challenge templates
	mathPool := []challengeTemplate{
		{
			activityType: "multiple_choice",
			questionJSON: mustJSON(map[string]interface{}{
				"prompt":      "🎯 Tantangan Harian: Berapakah hasil dari 23 + 48?",
				"options":     []string{"71", "61", "73", "68"},
				"answer":      "71",
				"hint":        "Tambahkan satuan dulu: 3+8=11, lalu puluhan: 20+40=60. Hasilnya 60+11!",
				"explanation": "23 + 48 = 71. Hebat!",
			}),
			maxScore: 10,
		},
		{
			activityType: "multiple_choice",
			questionJSON: mustJSON(map[string]interface{}{
				"prompt":      "🎯 Tantangan Harian: Berapakah hasil dari 9 × 7?",
				"options":     []string{"56", "63", "72", "54"},
				"answer":      "63",
				"hint":        "9 × 7 = (10 × 7) - 7 = 70 - 7",
				"explanation": "9 × 7 = 63. Kamu hebat!",
			}),
			maxScore: 10,
		},
		{
			activityType: "fill_blank",
			questionJSON: mustJSON(map[string]interface{}{
				"prompt":      "🎯 Tantangan Harian: Berapakah 100 - 37?",
				"answer":      "63",
				"hint":        "Kurangi puluhan: 100-30=70, lalu kurangi satuan: 70-7",
				"explanation": "100 - 37 = 63. Luar biasa!",
			}),
			maxScore: 10,
		},
		{
			activityType: "multiple_choice",
			questionJSON: mustJSON(map[string]interface{}{
				"prompt":      "🎯 Tantangan Harian: 156 + 89 = ?",
				"options":     []string{"245", "235", "255", "225"},
				"answer":      "245",
				"hint":        "150+80=230, lalu 6+9=15, jadi 230+15",
				"explanation": "156 + 89 = 245. Keren!",
			}),
			maxScore: 10,
		},
		{
			activityType: "fill_blank",
			questionJSON: mustJSON(map[string]interface{}{
				"prompt":      "🎯 Tantangan Harian: Berapakah 8 × 6?",
				"answer":      "48",
				"hint":        "8 × 6 = (8 × 5) + 8 = 40 + 8",
				"explanation": "8 × 6 = 48. Mantap!",
			}),
			maxScore: 10,
		},
	}

	sciencePool := []challengeTemplate{
		{
			activityType: "multiple_choice",
			questionJSON: mustJSON(map[string]interface{}{
				"prompt":      "🎯 Tantangan Harian: Planet manakah yang paling dekat dengan Matahari?",
				"options":     []string{"Merkurius", "Venus", "Bumi", "Mars"},
				"answer":      "Merkurius",
				"hint":        "Planet ini dinamai dewa pembawa pesan dari mitologi Romawi",
				"explanation": "Merkurius adalah planet terdekat dengan Matahari, berjarak sekitar 58 juta km.",
			}),
			maxScore: 10,
		},
		{
			activityType: "multiple_choice",
			questionJSON: mustJSON(map[string]interface{}{
				"prompt":      "🎯 Tantangan Harian: Air akan membeku pada suhu berapa derajat Celsius?",
				"options":     []string{"0°C", "10°C", "-10°C", "100°C"},
				"answer":      "0°C",
				"hint":        "Suhu ini adalah titik beku air",
				"explanation": "Air membeku pada suhu 0°C (32°F) pada tekanan standar.",
			}),
			maxScore: 10,
		},
	}

	languagePool := []challengeTemplate{
		{
			activityType: "fill_blank",
			questionJSON: mustJSON(map[string]interface{}{
				"prompt":      "🎯 Tantangan Harian: Lawan kata dari 'besar' adalah ...",
				"answer":      "kecil",
				"hint":        "Pikirkan benda yang ukurannya sangat mini",
				"explanation": "Lawan kata (antonim) dari 'besar' adalah 'kecil'.",
			}),
			maxScore: 10,
		},
		{
			activityType: "multiple_choice",
			questionJSON: mustJSON(map[string]interface{}{
				"prompt":      "🎯 Tantangan Harian: Manakah kata kerja (verba) di bawah ini?",
				"options":     []string{"Berlari", "Kucing", "Indah", "Cepat"},
				"answer":      "Berlari",
				"hint":        "Kata kerja menunjukkan tindakan atau perbuatan",
				"explanation": "'Berlari' adalah kata kerja karena menunjukkan aktivitas/tindakan.",
			}),
			maxScore: 10,
		},
	}

	// Combine all pools and pick randomly
	allTemplates := append(mathPool, sciencePool...)
	allTemplates = append(allTemplates, languagePool...)

	picked := allTemplates[r.Intn(len(allTemplates))]

	res, err := s.db.Exec(
		`INSERT INTO daily_challenges (challenge_date, age_group, activity_type, question_json, max_score)
		 VALUES (?, ?, ?, ?, ?)`,
		date, ageGroup, picked.activityType, picked.questionJSON, picked.maxScore,
	)
	if err != nil {
		return nil, fmt.Errorf("insert daily challenge: %w", err)
	}

	id, _ := res.LastInsertId()
	return &model.DailyChallenge{
		ID:            id,
		ChallengeDate: date,
		AgeGroup:      ageGroup,
		ActivityType:  picked.activityType,
		QuestionJSON:  picked.questionJSON,
		MaxScore:      picked.maxScore,
	}, nil
}

func mustJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ezedu/backend/internal/auth"
	"github.com/ezedu/backend/internal/store"
	"github.com/go-chi/chi/v5"
)

// LessonHandler handles endpoints for lessons, activities, and quiz submissions.
type LessonHandler struct {
	lessons    *store.LessonStore
	categories *store.CategoryStore
	progress   *store.ProgressStore
	children   *store.ChildStore
	badges     *store.BadgeStore
}

func NewLessonHandler(
	lessons *store.LessonStore,
	categories *store.CategoryStore,
	progress *store.ProgressStore,
	children *store.ChildStore,
	badges *store.BadgeStore,
) *LessonHandler {
	return &LessonHandler{
		lessons:    lessons,
		categories: categories,
		progress:   progress,
		children:   children,
		badges:     badges,
	}
}

// ListByCategory handles GET /api/categories/{slug}/lessons?age_group=builders&child_id=1
func (h *LessonHandler) ListByCategory(w http.ResponseWriter, r *http.Request) {
	categorySlug := chi.URLParam(r, "slug")
	ageGroup := r.URL.Query().Get("age_group")
	if ageGroup == "" {
		ageGroup = "builders" // Default
	}

	category, err := h.categories.GetBySlug(categorySlug)
	if err != nil || category == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Kategori tidak ditemukan"})
		return
	}

	lessons, err := h.lessons.ListByCategoryAndAgeGroup(category.ID, ageGroup)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal memuat pelajaran"})
		return
	}

	// Fetch progress if child_id provided
	childIDStr := r.URL.Query().Get("child_id")
	progressMap := make(map[int64]map[string]interface{})

	if childIDStr != "" {
		accountID := auth.AccountIDFromContext(r.Context())
		childID, _ := strconv.ParseInt(childIDStr, 10, 64)
		if childID > 0 {
			// Verify ownership
			child, err := h.children.GetByID(childID, accountID)
			if err == nil && child != nil {
				progList, err := h.progress.ListChildProgressByChild(childID)
				if err == nil {
					for _, p := range progList {
						progressMap[p.LessonID] = map[string]interface{}{
							"status":       p.Status,
							"score":        p.Score,
							"max_possible": p.MaxPossible,
							"attempts":     p.Attempts,
						}
					}
				}
			}
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"category": category,
		"lessons":  lessons,
		"progress": progressMap,
	})
}

// GetByID handles GET /api/lessons/{id}
func (h *LessonHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	lessonID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID pelajaran tidak valid"})
		return
	}

	lesson, err := h.lessons.GetByID(lessonID)
	if err != nil || lesson == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Pelajaran tidak ditemukan"})
		return
	}

	activities, err := h.lessons.ListActivitiesByLessonID(lessonID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal memuat soal"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"lesson":     lesson,
		"activities": activities,
	})
}

type submitActivityReq struct {
	ChildID   int64           `json:"child_id"`
	Answer    json.RawMessage `json:"answer"`
	AttemptNo int             `json:"attempt_no"`
}

// SubmitActivity handles POST /api/activities/{id}/submit
func (h *LessonHandler) SubmitActivity(w http.ResponseWriter, r *http.Request) {
	activityID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID soal tidak valid"})
		return
	}

	var req submitActivityReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format jawaban tidak valid"})
		return
	}

	activity, err := h.lessons.GetActivityByID(activityID)
	if err != nil || activity == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Soal tidak ditemukan"})
		return
	}

	// Parse question details
	var qData struct {
		Answer        string   `json:"answer"`
		ExpectedOrder []string `json:"expected_order"`
		Hint          string   `json:"hint"`
		Explanation   string   `json:"explanation"`
	}
	_ = json.Unmarshal([]byte(activity.QuestionJSON), &qData)

	isCorrect := false
	score := 0

	// Evaluate answer based on activity type
	switch activity.Type {
	case "multiple_choice", "fill_blank":
		var userChoice string
		_ = json.Unmarshal(req.Answer, &userChoice)
		if strings.EqualFold(strings.TrimSpace(userChoice), strings.TrimSpace(qData.Answer)) {
			isCorrect = true
			score = activity.MaxScore
		}

	case "drag_drop", "sequencing", "block_code":
		var userOrder []string
		_ = json.Unmarshal(req.Answer, &userOrder)

		if len(userOrder) == len(qData.ExpectedOrder) {
			match := true
			for i, v := range qData.ExpectedOrder {
				if i >= len(userOrder) {
					match = false
					break
				}
				u := normalizeString(userOrder[i])
				exp := normalizeString(v)
				if u != exp && !strings.Contains(u, exp) && !strings.Contains(exp, u) {
					match = false
					break
				}
			}
			if match {
				isCorrect = true
				score = activity.MaxScore
			}
		}

	case "drawing", "pixel_art":
		isCorrect = true
		score = activity.MaxScore

	case "timed":
		// Timed activities (e.g., Math Racer) send their score from the client
		var clientScore int
		_ = json.Unmarshal(req.Answer, &clientScore)
		if clientScore > 0 {
			isCorrect = true
			score = clientScore
			if score > activity.MaxScore {
				score = activity.MaxScore
			}
		}
	}

	if req.AttemptNo == 0 {
		req.AttemptNo = 1
	}

	// Record activity result if childID provided
	if req.ChildID > 0 {
		_ = h.progress.RecordActivityResult(req.ChildID, activityID, string(req.Answer), isCorrect, score, req.AttemptNo)
	}

	feedbackMsg := "Coba lagi! Kamu pasti bisa 💪"
	if isCorrect {
		feedbackMsg = "Hebat sekali! Jawabanmu benar 🎉"
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"is_correct":  isCorrect,
		"score":       score,
		"max_score":   activity.MaxScore,
		"feedback":    feedbackMsg,
		"hint":        qData.Hint,
		"explanation": qData.Explanation,
	})
}

type completeLessonReq struct {
	ChildID      int64 `json:"child_id"`
	FinalScore   int   `json:"final_score"`
	MaxScore     int   `json:"max_score"`
	TimeSpentSec int   `json:"time_spent_sec"`
}

// CompleteLesson handles POST /api/lessons/{id}/complete
func (h *LessonHandler) CompleteLesson(w http.ResponseWriter, r *http.Request) {
	lessonID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID pelajaran tidak valid"})
		return
	}

	var req completeLessonReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
		return
	}

	accountID := auth.AccountIDFromContext(r.Context())
	child, err := h.children.GetByID(req.ChildID, accountID)
	if err != nil || child == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Profil anak tidak valid"})
		return
	}

	lesson, err := h.lessons.GetByID(lessonID)
	if err != nil || lesson == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Pelajaran tidak ditemukan"})
		return
	}

	xpReward := lesson.XPReward
	if req.FinalScore < (req.MaxScore / 2) {
		xpReward = xpReward / 2 // Half reward for lower score
	}

	if err := h.progress.CompleteLesson(req.ChildID, lessonID, req.FinalScore, req.MaxScore, req.TimeSpentSec, xpReward); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal menyimpan progres pelajaran"})
		return
	}

	// Update streak
	_ = h.progress.UpdateStreak(req.ChildID)

	// Evaluate and award badges
	var awardedBadges []map[string]interface{}
	if h.badges != nil {
		newBadges, _ := h.badges.EvaluateAndAwardBadges(req.ChildID)
		for _, b := range newBadges {
			awardedBadges = append(awardedBadges, map[string]interface{}{
				"name":        b.Name,
				"description": b.Description,
				"icon":        b.Icon,
				"slug":        b.Slug,
			})
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message":        "Selamat! Kamu berhasil menyelesaikan pelajaran 🎉",
		"xp_earned":      xpReward,
		"badges_awarded": awardedBadges,
	})
}

// GetChildProgress handles GET /api/children/{id}/progress
func (h *LessonHandler) GetChildProgress(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	childID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID tidak valid"})
		return
	}

	child, err := h.children.GetByID(childID, accountID)
	if err != nil || child == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Profil anak tidak ditemukan"})
		return
	}

	summary, err := h.progress.GetChildProgressSummary(childID, child.AgeGroup)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal memuat data kemajuan"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"child":    child,
		"progress": summary,
	})
}

// GetChildBadges handles GET /api/children/{id}/badges
func (h *LessonHandler) GetChildBadges(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	childID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID tidak valid"})
		return
	}

	child, err := h.children.GetByID(childID, accountID)
	if err != nil || child == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Profil anak tidak ditemukan"})
		return
	}

	badges, err := h.badges.ListChildBadges(childID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal memuat data lencana"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"badges": badges,
	})
}

func normalizeString(s string) string {
	s = strings.ToLower(s)
	var sb strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == ' ' {
			sb.WriteRune(r)
		}
	}
	return strings.TrimSpace(sb.String())
}


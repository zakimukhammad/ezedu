package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ezedu/backend/internal/store"
)

// DailyHandler handles daily challenge API endpoints.
type DailyHandler struct {
	daily *store.DailyStore
}

func NewDailyHandler(daily *store.DailyStore) *DailyHandler {
	return &DailyHandler{daily: daily}
}

// GetToday handles GET /api/daily-challenge?age_group=builders&child_id=1
func (h *DailyHandler) GetToday(w http.ResponseWriter, r *http.Request) {
	ageGroup := r.URL.Query().Get("age_group")
	if ageGroup == "" {
		ageGroup = "builders"
	}

	challenge, err := h.daily.GetTodayChallenge(ageGroup)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal memuat tantangan harian"})
		return
	}

	// Check if child has completed
	childIDStr := r.URL.Query().Get("child_id")
	completed := false
	if childIDStr != "" {
		childID, _ := strconv.ParseInt(childIDStr, 10, 64)
		if childID > 0 {
			completed = h.daily.HasCompleted(childID, challenge.ID)
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"challenge": challenge,
		"completed": completed,
	})
}

// Submit handles POST /api/daily-challenge/submit
func (h *DailyHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ChildID     int64           `json:"child_id"`
		ChallengeID int64           `json:"challenge_id"`
		Answer      json.RawMessage `json:"answer"`
		Score       int             `json:"score"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format data tidak valid"})
		return
	}

	if req.ChildID == 0 || req.ChallengeID == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "child_id dan challenge_id wajib diisi"})
		return
	}

	// Check if already completed
	if h.daily.HasCompleted(req.ChildID, req.ChallengeID) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"already_completed": true,
			"message":           "Kamu sudah menyelesaikan tantangan hari ini! 🎉",
		})
		return
	}

	if err := h.daily.SubmitResult(req.ChildID, req.ChallengeID, req.Score); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal menyimpan hasil"})
		return
	}

	streak := h.daily.GetStreak(req.ChildID)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"score":   req.Score,
		"streak":  streak,
		"message": "Tantangan harian selesai! Hebat! 🏆",
	})
}

// GetStreak handles GET /api/daily-challenge/streak?child_id=1
func (h *DailyHandler) GetStreak(w http.ResponseWriter, r *http.Request) {
	childIDStr := r.URL.Query().Get("child_id")
	if childIDStr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "child_id wajib diisi"})
		return
	}

	childID, _ := strconv.ParseInt(childIDStr, 10, 64)
	streak := h.daily.GetStreak(childID)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"child_id": childID,
		"streak":   streak,
	})
}

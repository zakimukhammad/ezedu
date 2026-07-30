package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ezedu/backend/internal/auth"
	"github.com/ezedu/backend/internal/model"
	"github.com/ezedu/backend/internal/store"
	"github.com/go-chi/chi/v5"
)

// LeaderboardHandler handles leaderboard endpoints.
type LeaderboardHandler struct {
	leaderboard *store.LeaderboardStore
	children    *store.ChildStore
}

func NewLeaderboardHandler(leaderboard *store.LeaderboardStore, children *store.ChildStore) *LeaderboardHandler {
	return &LeaderboardHandler{
		leaderboard: leaderboard,
		children:    children,
	}
}

type toggleOptInRequest struct {
	OptIn bool `json:"opt_in"`
}

// GetWeekly handles GET /api/leaderboard?child_id={id}
func (h *LeaderboardHandler) GetWeekly(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	childIDStr := r.URL.Query().Get("child_id")
	if childIDStr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "child_id wajib diisi"})
		return
	}

	childID, err := strconv.ParseInt(childIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "child_id tidak valid"})
		return
	}

	// Verify child profile exists and belongs to this account
	child, err := h.children.GetByID(childID, accountID)
	if err != nil || child == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Profil anak tidak ditemukan"})
		return
	}

	weekStart := store.CurrentWeekStart()

	// Get leaderboard top 20
	entries, err := h.leaderboard.GetWeeklyLeaderboard(childID, weekStart, 20)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal memuat papan peringkat"})
		return
	}

	if entries == nil {
		entries = []model.LeaderboardEntry{}
	}

	// Get child's current rank
	rank, _, _ := h.leaderboard.GetChildRankAndEntry(childID, weekStart)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"week_start":   weekStart,
		"opt_in":       child.LeaderboardOptIn,
		"child_rank":   rank,
		"is_challenger": child.AgeGroup == "challengers",
		"entries":      entries,
	})
}

// ToggleOptIn handles PUT /api/children/{id}/leaderboard-opt-in
func (h *LeaderboardHandler) ToggleOptIn(w http.ResponseWriter, r *http.Request) {
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

	var req toggleOptInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format data tidak valid"})
		return
	}

	if err := h.children.UpdateLeaderboardOptIn(childID, accountID, req.OptIn); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal memperbarui opsi papan peringkat"})
		return
	}

	weekStart := store.CurrentWeekStart()
	if req.OptIn {
		// If opting in, insert initial entry with 0 gained if not present
		_ = h.leaderboard.UpsertWeeklyXP(childID, weekStart, 0, child.AvatarID)
	} else {
		// If opting out, remove current week entry
		_ = h.leaderboard.RemoveChildEntry(childID, weekStart)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Pengaturan papan peringkat berhasil diperbarui!",
		"opt_in":  req.OptIn,
	})
}

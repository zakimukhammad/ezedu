package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ezedu/backend/internal/auth"
	"github.com/ezedu/backend/internal/store"
	"github.com/go-chi/chi/v5"
)

// ChildHandler handles child profile endpoints.
type ChildHandler struct {
	children *store.ChildStore
	progress *store.ProgressStore
}

func NewChildHandler(children *store.ChildStore, progress *store.ProgressStore) *ChildHandler {
	return &ChildHandler{children: children, progress: progress}
}

type createChildRequest struct {
	Name      string `json:"name"`
	BirthYear int    `json:"birth_year"`
	AvatarID  int    `json:"avatar_id"`
}

type updateChildRequest struct {
	Name      string `json:"name"`
	BirthYear int    `json:"birth_year"`
	AvatarID  int    `json:"avatar_id"`
}

// List handles GET /api/children
func (h *ChildHandler) List(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	children, err := h.children.ListByAccount(accountID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal memuat profil anak"})
		return
	}
	if children == nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{"children": []interface{}{}})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"children": children})
}

// Create handles POST /api/children
func (h *ChildHandler) Create(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())

	var req createChildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
		return
	}

	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Nama anak wajib diisi"})
		return
	}
	currentYear := time.Now().Year()
	if req.BirthYear < currentYear-16 || req.BirthYear > currentYear {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Tahun lahir tidak valid"})
		return
	}

	// Check max 4 children per account
	count, err := h.children.CountByAccount(accountID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal memeriksa jumlah profil"})
		return
	}
	if count >= 4 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Maksimal 4 profil anak per akun"})
		return
	}

	if req.AvatarID == 0 {
		req.AvatarID = 1
	}

	child, err := h.children.Create(accountID, req.Name, req.BirthYear, req.AvatarID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal membuat profil anak"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Profil anak berhasil dibuat!",
		"child":   child,
	})
}

// Update handles PUT /api/children/{id}
func (h *ChildHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	childID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID tidak valid"})
		return
	}

	// Check ownership
	existing, err := h.children.GetByID(childID, accountID)
	if err != nil || existing == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Profil anak tidak ditemukan"})
		return
	}

	var req updateChildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
		return
	}

	if req.Name == "" {
		req.Name = existing.Name
	}
	if req.BirthYear == 0 {
		req.BirthYear = existing.BirthYear
	}
	if req.AvatarID == 0 {
		req.AvatarID = existing.AvatarID
	}

	currentYear := time.Now().Year()
	if req.BirthYear < currentYear-16 || req.BirthYear > currentYear {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Tahun lahir tidak valid"})
		return
	}

	if err := h.children.Update(childID, accountID, req.Name, req.BirthYear, req.AvatarID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal memperbarui profil"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Profil berhasil diperbarui!"})
}

// Delete handles DELETE /api/children/{id}
func (h *ChildHandler) Delete(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	childID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID tidak valid"})
		return
	}

	existing, err := h.children.GetByID(childID, accountID)
	if err != nil || existing == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Profil anak tidak ditemukan"})
		return
	}

	if err := h.children.Delete(childID, accountID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal menghapus profil"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Profil berhasil dihapus"})
}

type updateDailyLimitRequest struct {
	DailyLimitMin *int `json:"daily_limit_min"` // nil or int > 0
}

// UpdateDailyLimit handles PUT /api/children/{id}/daily-limit
func (h *ChildHandler) UpdateDailyLimit(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	childID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID tidak valid"})
		return
	}

	existing, err := h.children.GetByID(childID, accountID)
	if err != nil || existing == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Profil anak tidak ditemukan"})
		return
	}

	var req updateDailyLimitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format data tidak valid"})
		return
	}

	if req.DailyLimitMin != nil {
		if *req.DailyLimitMin <= 0 || *req.DailyLimitMin > 480 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Batas waktu harus antara 1 dan 480 menit"})
			return
		}
	}

	if err := h.children.UpdateDailyLimit(childID, accountID, req.DailyLimitMin); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal memperbarui batas waktu"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message":         "Batas waktu harian berhasil diperbarui",
		"daily_limit_min": req.DailyLimitMin,
	})
}

// GetRemainingTime handles GET /api/children/{id}/remaining-time
func (h *ChildHandler) GetRemainingTime(w http.ResponseWriter, r *http.Request) {
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

	if child.DailyLimitMin == nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"has_limit": false,
		})
		return
	}

	limitMin := *child.DailyLimitMin
	limitSec := limitMin * 60

	todayTimeSpentSec, err := h.progress.GetTodayTimeSpent(childID)
	if err != nil {
		todayTimeSpentSec = 0
	}

	remainingSec := limitSec - todayTimeSpentSec

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"has_limit":         true,
		"daily_limit_min":   limitMin,
		"time_used_sec":     todayTimeSpentSec,
		"remaining_sec":     remainingSec,
		"limit_reached":     remainingSec <= 0,
	})
}


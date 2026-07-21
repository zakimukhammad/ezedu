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
}

func NewChildHandler(children *store.ChildStore) *ChildHandler {
	return &ChildHandler{children: children}
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

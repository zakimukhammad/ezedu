package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ezedu/backend/internal/auth"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	service *auth.Service
}

func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

// Signup handles POST /api/auth/signup
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req auth.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
		return
	}

	account, session, err := h.service.Signup(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	auth.SetSessionCookie(w, session)
	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Akun berhasil dibuat!",
		"account": map[string]interface{}{
			"id":          account.ID,
			"email":       account.Email,
			"parent_name": account.ParentName,
		},
	})
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
		return
	}

	ipAddress := r.RemoteAddr
	account, session, err := h.service.Login(req, ipAddress)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	auth.SetSessionCookie(w, session)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Berhasil masuk!",
		"account": map[string]interface{}{
			"id":          account.ID,
			"email":       account.Email,
			"parent_name": account.ParentName,
		},
	})
}

// Logout handles POST /api/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := auth.GetSessionToken(r)
	if token != "" {
		h.service.Logout(token)
	}
	auth.ClearSessionCookie(w)
	writeJSON(w, http.StatusOK, map[string]string{"message": "Berhasil keluar"})
}

// Me handles GET /api/auth/me — returns the current authenticated account.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	if accountID == 0 {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Tidak terautentikasi"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"account_id": accountID,
	})
}

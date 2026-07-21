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

type setPINRequest struct {
	PIN string `json:"pin"`
}

type verifyPINRequest struct {
	PIN string `json:"pin"`
}

// UpdatePIN handles PUT /api/auth/pin
func (h *AuthHandler) UpdatePIN(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	if accountID == 0 {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Tidak terautentikasi"})
		return
	}

	var req setPINRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format data tidak valid"})
		return
	}

	if err := h.service.SetPIN(accountID, req.PIN); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "PIN Orang Tua berhasil diperbarui!"})
}

// VerifyPIN handles POST /api/auth/pin/verify
func (h *AuthHandler) VerifyPIN(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	if accountID == 0 {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Tidak terautentikasi"})
		return
	}

	var req verifyPINRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format data tidak valid"})
		return
	}

	valid, err := h.service.VerifyPIN(accountID, req.PIN)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal memverifikasi PIN"})
		return
	}

	if !valid {
		writeJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"valid": false,
			"error": "PIN Orang Tua salah. Silakan coba lagi.",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"valid":   true,
		"message": "PIN benar",
	})
}

// Me handles GET /api/auth/me — returns the current authenticated account.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	if accountID == 0 {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Tidak terautentikasi"})
		return
	}

	account, err := h.service.GetAccountByID(accountID)
	if err != nil || account == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Akun tidak ditemukan"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"account_id":     account.ID,
		"email":          account.Email,
		"parent_name":    account.ParentName,
		"has_parent_pin": account.ParentPIN != "",
	})
}


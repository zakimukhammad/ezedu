package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ezedu/backend/internal/model"
	"github.com/ezedu/backend/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const (
	accountContextKey contextKey = "account"
	sessionCookieName string     = "ezedu_session"
)

// Service handles authentication logic.
type Service struct {
	accounts *store.AccountStore
	sessions *store.SessionStore
}

func NewService(accounts *store.AccountStore, sessions *store.SessionStore) *Service {
	return &Service{accounts: accounts, sessions: sessions}
}

// SignupRequest represents the signup payload.
type SignupRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	ParentName string `json:"parent_name"`
}

// LoginRequest represents the login payload.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Signup creates a new account.
func (s *Service) Signup(req SignupRequest) (*model.Account, *model.Session, error) {
	// Validate
	if strings.TrimSpace(req.Email) == "" {
		return nil, nil, fmt.Errorf("email wajib diisi")
	}
	if len(req.Password) < 6 {
		return nil, nil, fmt.Errorf("kata sandi minimal 6 karakter")
	}
	if strings.TrimSpace(req.ParentName) == "" {
		return nil, nil, fmt.Errorf("nama orang tua wajib diisi")
	}

	// Check existing
	existing, err := s.accounts.GetByEmail(req.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal memeriksa email: %w", err)
	}
	if existing != nil {
		return nil, nil, fmt.Errorf("email sudah terdaftar")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal mengenkripsi kata sandi: %w", err)
	}

	// Create account
	id, err := s.accounts.Create(req.Email, string(hash), req.ParentName)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal membuat akun: %w", err)
	}

	account, err := s.accounts.GetByID(id)
	if err != nil {
		return nil, nil, err
	}

	// Create session
	session, err := s.sessions.Create(id, "")
	if err != nil {
		return nil, nil, err
	}

	return account, session, nil
}

// Login authenticates a user.
func (s *Service) Login(req LoginRequest, ipAddress string) (*model.Account, *model.Session, error) {
	if strings.TrimSpace(req.Email) == "" {
		return nil, nil, fmt.Errorf("email wajib diisi")
	}
	if strings.TrimSpace(req.Password) == "" {
		return nil, nil, fmt.Errorf("kata sandi wajib diisi")
	}

	account, err := s.accounts.GetByEmail(req.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal memeriksa akun: %w", err)
	}
	if account == nil {
		return nil, nil, fmt.Errorf("email atau kata sandi salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(req.Password)); err != nil {
		return nil, nil, fmt.Errorf("email atau kata sandi salah")
	}

	session, err := s.sessions.Create(account.ID, ipAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal membuat sesi: %w", err)
	}

	return account, session, nil
}

// Logout destroys a session.
func (s *Service) Logout(token string) error {
	return s.sessions.Delete(token)
}

// GetAccountByID retrieves account details by ID.
func (s *Service) GetAccountByID(id int64) (*model.Account, error) {
	return s.accounts.GetByID(id)
}

// SetPIN sets or updates a 4-digit Parent PIN.
func (s *Service) SetPIN(accountID int64, pin string) error {
	pin = strings.TrimSpace(pin)
	if len(pin) != 4 {
		return fmt.Errorf("PIN orang tua harus 4 digit angka")
	}
	for _, ch := range pin {
		if ch < '0' || ch > '9' {
			return fmt.Errorf("PIN orang tua harus berupa angka")
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gagal mengenkripsi PIN: %w", err)
	}

	return s.accounts.UpdatePIN(accountID, string(hash))
}

// VerifyPIN verifies the provided 4-digit Parent PIN against the stored hash.
func (s *Service) VerifyPIN(accountID int64, pin string) (bool, error) {
	account, err := s.accounts.GetByID(accountID)
	if err != nil || account == nil {
		return false, fmt.Errorf("akun tidak ditemukan")
	}

	// If no PIN has been set yet, verification passes
	if account.ParentPIN == "" {
		return true, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.ParentPIN), []byte(pin))
	if err != nil {
		return false, nil
	}
	return true, nil
}


// SetSessionCookie sets the session cookie on the response.
func SetSessionCookie(w http.ResponseWriter, session *model.Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    session.ID,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // Set to true in production with HTTPS
	})
}

// ClearSessionCookie removes the session cookie.
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

// SessionMiddleware validates the session cookie and adds account to context.
func SessionMiddleware(sessions *store.SessionStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionCookieName)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "Sesi tidak ditemukan. Silakan masuk kembali.")
				return
			}

			session, err := sessions.GetByToken(cookie.Value)
			if err != nil || session == nil {
				ClearSessionCookie(w)
				writeError(w, http.StatusUnauthorized, "Sesi sudah berakhir. Silakan masuk kembali.")
				return
			}

			ctx := context.WithValue(r.Context(), accountContextKey, session.AccountID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AccountIDFromContext extracts the account ID from the request context.
func AccountIDFromContext(ctx context.Context) int64 {
	id, ok := ctx.Value(accountContextKey).(int64)
	if !ok {
		return 0
	}
	return id
}

// GetSessionToken extracts the session token from the cookie.
func GetSessionToken(r *http.Request) string {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

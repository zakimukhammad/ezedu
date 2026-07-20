package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ezedu/backend/internal/auth"
	"github.com/ezedu/backend/internal/handler"
	"github.com/ezedu/backend/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Determine data directory
	dataDir := os.Getenv("EZEDU_DATA_DIR")
	if dataDir == "" {
		dataDir = filepath.Join(".", "data")
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	dbPath := filepath.Join(dataDir, "ezedu.db")
	log.Printf("Database path: %s", dbPath)

	// Initialize SQLite store
	db, err := store.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := store.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations complete")

	// Seed default data
	if err := store.SeedDefaults(db); err != nil {
		log.Fatalf("Failed to seed defaults: %v", err)
	}

	// Initialize stores
	accountStore := store.NewAccountStore(db)
	sessionStore := store.NewSessionStore(db)
	childStore := store.NewChildStore(db)
	categoryStore := store.NewCategoryStore(db)

	// Initialize auth service
	authService := auth.NewService(accountStore, sessionStore)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	childHandler := handler.NewChildHandler(childStore)
	categoryHandler := handler.NewCategoryHandler(categoryStore)

	// Build router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(corsMiddleware)

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Public auth routes
		r.Post("/auth/signup", authHandler.Signup)
		r.Post("/auth/login", authHandler.Login)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(auth.SessionMiddleware(sessionStore))

			r.Post("/auth/logout", authHandler.Logout)
			r.Get("/auth/me", authHandler.Me)

			// Child profiles
			r.Get("/children", childHandler.List)
			r.Post("/children", childHandler.Create)
			r.Put("/children/{id}", childHandler.Update)
			r.Delete("/children/{id}", childHandler.Delete)

			// Categories
			r.Get("/categories", categoryHandler.List)
		})
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Start server
	port := os.Getenv("EZEDU_PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("EzEdu backend starting on %s", addr)

	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

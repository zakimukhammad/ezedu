package handler

import (
	"net/http"

	"github.com/ezedu/backend/internal/store"
)

// CategoryHandler handles category endpoints.
type CategoryHandler struct {
	categories *store.CategoryStore
}

func NewCategoryHandler(categories *store.CategoryStore) *CategoryHandler {
	return &CategoryHandler{categories: categories}
}

// List handles GET /api/categories
func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categories.List()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal memuat kategori"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"categories": categories})
}

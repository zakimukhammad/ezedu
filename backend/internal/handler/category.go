package handler

import (
	"net/http"

	"github.com/ezedu/backend/internal/model"
	"github.com/ezedu/backend/internal/store"
)

// CategoryHandler handles category endpoints.
type CategoryHandler struct {
	categories *store.CategoryStore
}

func NewCategoryHandler(categories *store.CategoryStore) *CategoryHandler {
	return &CategoryHandler{categories: categories}
}

// List handles GET /api/categories?age_group=toddlers
func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categories.List()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal memuat kategori"})
		return
	}

	ageGroup := r.URL.Query().Get("age_group")
	if ageGroup != "" {
		filtered := make([]model.Category, 0)
		for _, c := range categories {
			if ageGroup == "toddlers" {
				if c.Slug == "toddlers" {
					filtered = append(filtered, c)
				}
			} else {
				if c.Slug != "toddlers" {
					filtered = append(filtered, c)
				}
			}
		}
		categories = filtered
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"categories": categories})
}

package store

import (
	"database/sql"
	"fmt"

	"github.com/ezedu/backend/internal/model"
)

// CategoryStore handles category persistence.
type CategoryStore struct {
	db *sql.DB
}

func NewCategoryStore(db *sql.DB) *CategoryStore {
	return &CategoryStore{db: db}
}

// List retrieves all categories in sort order.
func (s *CategoryStore) List() ([]model.Category, error) {
	rows, err := s.db.Query(
		`SELECT id, slug, name, description, icon, color, sort_order
		 FROM categories ORDER BY sort_order`,
	)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Slug, &c.Name, &c.Description, &c.Icon, &c.Color, &c.SortOrder); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// GetBySlug retrieves a category by slug.
func (s *CategoryStore) GetBySlug(slug string) (*model.Category, error) {
	c := &model.Category{}
	err := s.db.QueryRow(
		`SELECT id, slug, name, description, icon, color, sort_order
		 FROM categories WHERE slug = ?`, slug,
	).Scan(&c.ID, &c.Slug, &c.Name, &c.Description, &c.Icon, &c.Color, &c.SortOrder)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get category: %w", err)
	}
	return c, nil
}

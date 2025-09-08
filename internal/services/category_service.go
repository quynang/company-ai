package services

import (
	"company-ai-training/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{
		db: db,
	}
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(name, description string) (*models.Category, error) {
	category := &models.Category{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategory retrieves a category by ID
func (s *CategoryService) GetCategory(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	if err := s.db.First(&category, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// GetCategoryByName retrieves a category by name
func (s *CategoryService) GetCategoryByName(name string) (*models.Category, error) {
	var category models.Category
	if err := s.db.First(&category, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// GetAllCategories retrieves all categories
func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := s.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// UpdateCategory updates a category
func (s *CategoryService) UpdateCategory(id uuid.UUID, name, description string) (*models.Category, error) {
	category := &models.Category{
		ID:          id,
		Name:        name,
		Description: description,
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Model(&models.Category{}).Where("id = ?", id).Updates(category).Error; err != nil {
		return nil, err
	}

	// Return updated category
	return s.GetCategory(id)
}

// DeleteCategory deletes a category
func (s *CategoryService) DeleteCategory(id uuid.UUID) error {
	// Check if category is being used by any documents
	var count int64
	if err := s.db.Model(&models.Document{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("cannot delete category: it is being used by documents")
	}

	// Check if category is being used by any chat sessions
	if err := s.db.Model(&models.ChatSession{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("cannot delete category: it is being used by chat sessions")
	}

	// Delete category
	return s.db.Delete(&models.Category{}, "id = ?", id).Error
}

// GetDocumentsByCategory retrieves all documents in a specific category
func (s *CategoryService) GetDocumentsByCategory(categoryID uuid.UUID) ([]models.Document, error) {
	var documents []models.Document
	if err := s.db.Where("category_id = ?", categoryID).Find(&documents).Error; err != nil {
		return nil, err
	}

	// Populate chunks count for each document
	for i := range documents {
		var count int64
		if err := s.db.Model(&models.DocumentChunk{}).Where("document_id = ?", documents[i].ID).Count(&count).Error; err != nil {
			continue // Skip if error, keep count as 0
		}
		documents[i].ChunksCount = int(count)
	}

	return documents, nil
}


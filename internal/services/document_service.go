package services

import (
	"bytes"
	"company-ai-training/internal/models"
	"errors"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ledongthuc/pdf"
	"github.com/unidoc/unioffice/document"
	"gorm.io/gorm"
)

type DocumentService struct {
	db *gorm.DB
}

func NewDocumentService(db *gorm.DB) *DocumentService {
	return &DocumentService{
		db: db,
	}
}

func (s *DocumentService) UploadDocument(file *multipart.FileHeader) (*models.Document, error) {
	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" && ext != ".docx" && ext != ".txt" {
		return nil, errors.New("unsupported file type. Only PDF, DOCX, and TXT files are allowed")
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Extract content based on file type
	var content string
	switch ext {
	case ".pdf":
		content, err = s.extractPDFContent(src)
	case ".docx":
		content, err = s.extractDocxContent(src)
	case ".txt":
		content, err = s.extractTextContent(src)
	}

	if err != nil {
		return nil, err
	}

	// Create document record
	doc := &models.Document{
		ID:         uuid.New(),
		Name:       file.Filename,
		Content:    content,
		Type:       strings.TrimPrefix(ext, "."),
		Size:       int64(len(content)),
		UploadedAt: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Save to database
	if err := s.db.Create(doc).Error; err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *DocumentService) CreateDocumentFromText(name, content string) (*models.Document, error) {
	// Create document record from text content
	doc := &models.Document{
		ID:         uuid.New(),
		Name:       name,
		Content:    content,
		Type:       "txt",
		Size:       int64(len(content)),
		UploadedAt: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Save to database
	if err := s.db.Create(doc).Error; err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *DocumentService) CreateDocumentFromTextWithCategories(name, content string, categoryIDs []uuid.UUID) (*models.Document, error) {
	// Create document record from text content with categories
	doc := &models.Document{
		ID:         uuid.New(),
		Name:       name,
		Content:    content,
		Type:       "txt",
		Size:       int64(len(content)),
		UploadedAt: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Save to database
	if err := s.db.Create(doc).Error; err != nil {
		return nil, err
	}

	// Associate categories if provided
	if len(categoryIDs) > 0 {
		if err := s.AssignCategoriesToDocument(doc.ID, categoryIDs); err != nil {
			return nil, err
		}
	}

	return doc, nil
}

func (s *DocumentService) UploadDocumentWithCategories(file *multipart.FileHeader, categoryIDs []uuid.UUID) (*models.Document, error) {
	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" && ext != ".docx" && ext != ".txt" {
		return nil, errors.New("unsupported file type. Only PDF, DOCX, and TXT files are allowed")
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Extract content based on file type
	var content string
	switch ext {
	case ".pdf":
		content, err = s.extractPDFContent(src)
	case ".docx":
		content, err = s.extractDocxContent(src)
	case ".txt":
		content, err = s.extractTextContent(src)
	}

	if err != nil {
		return nil, err
	}

	// Create document record with categories
	doc := &models.Document{
		ID:         uuid.New(),
		Name:       file.Filename,
		Content:    content,
		Type:       strings.TrimPrefix(ext, "."),
		Size:       int64(len(content)),
		UploadedAt: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Save to database
	if err := s.db.Create(doc).Error; err != nil {
		return nil, err
	}

	// Associate categories if provided
	if len(categoryIDs) > 0 {
		if err := s.AssignCategoriesToDocument(doc.ID, categoryIDs); err != nil {
			return nil, err
		}
	}

	return doc, nil
}

func (s *DocumentService) GetDocument(id uuid.UUID) (*models.Document, error) {
	var doc models.Document
	if err := s.db.Preload("Categories").First(&doc, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Populate chunks count
	var count int64
	if err := s.db.Model(&models.DocumentChunk{}).Where("document_id = ?", doc.ID).Count(&count).Error; err != nil {
		doc.ChunksCount = 0
	} else {
		doc.ChunksCount = int(count)
	}

	return &doc, nil
}

func (s *DocumentService) GetAllDocuments() ([]models.Document, error) {
	var docs []models.Document
	if err := s.db.Preload("Categories").Find(&docs).Error; err != nil {
		return nil, err
	}

	// Populate chunks count for each document
	for i := range docs {
		var count int64
		if err := s.db.Model(&models.DocumentChunk{}).Where("document_id = ?", docs[i].ID).Count(&count).Error; err != nil {
			continue // Skip if error, keep count as 0
		}
		docs[i].ChunksCount = int(count)
	}

	return docs, nil
}

func (s *DocumentService) GetDocumentsByCategory(categoryID uuid.UUID) ([]models.Document, error) {
	var docs []models.Document
	if err := s.db.Preload("Categories").
		Joins("JOIN document_categories ON documents.id = document_categories.document_id").
		Where("document_categories.category_id = ?", categoryID).
		Find(&docs).Error; err != nil {
		return nil, err
	}

	// Populate chunks count for each document
	for i := range docs {
		var count int64
		if err := s.db.Model(&models.DocumentChunk{}).Where("document_id = ?", docs[i].ID).Count(&count).Error; err != nil {
			continue // Skip if error, keep count as 0
		}
		docs[i].ChunksCount = int(count)
	}

	return docs, nil
}

func (s *DocumentService) AssignCategoriesToDocument(documentID uuid.UUID, categoryIDs []uuid.UUID) error {
	// First, remove existing associations
	if err := s.db.Where("document_id = ?", documentID).Delete(&models.DocumentCategory{}).Error; err != nil {
		return err
	}

	// Then, create new associations
	for _, categoryID := range categoryIDs {
		docCategory := &models.DocumentCategory{
			DocumentID: documentID,
			CategoryID: categoryID,
			CreatedAt:  time.Now(),
		}
		if err := s.db.Create(docCategory).Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *DocumentService) GetDocumentCategories(documentID uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	if err := s.db.Joins("JOIN document_categories ON categories.id = document_categories.category_id").
		Where("document_categories.document_id = ?", documentID).
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *DocumentService) DeleteDocument(id uuid.UUID) error {
	// Delete document chunks first
	if err := s.db.Where("document_id = ?", id).Delete(&models.DocumentChunk{}).Error; err != nil {
		return err
	}

	// Delete document-category associations
	if err := s.db.Where("document_id = ?", id).Delete(&models.DocumentCategory{}).Error; err != nil {
		return err
	}

	// Delete document
	if err := s.db.Delete(&models.Document{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

func (s *DocumentService) UpdateDocument(id uuid.UUID, content string) error {
	// Delete existing chunks first
	if err := s.db.Where("document_id = ?", id).Delete(&models.DocumentChunk{}).Error; err != nil {
		return err
	}

	// Update document content and timestamp
	if err := s.db.Model(&models.Document{}).Where("id = ?", id).Updates(map[string]interface{}{
		"content":    content,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return err
	}

	return nil
}

func (s *DocumentService) extractPDFContent(reader io.Reader) (string, error) {
	// Read all data into memory
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	// Create a bytes reader that implements io.ReaderAt
	bytesReader := bytes.NewReader(data)

	pdfReader, err := pdf.NewReader(bytesReader, int64(len(data)))
	if err != nil {
		return "", err
	}

	var content strings.Builder
	numPages := pdfReader.NumPage()

	for i := 1; i <= numPages; i++ {
		page := pdfReader.Page(i)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			continue // Skip pages with extraction errors
		}

		content.WriteString(text)
		content.WriteString("\n")
	}

	return content.String(), nil
}

func (s *DocumentService) extractDocxContent(reader io.Reader) (string, error) {
	// Read all data into memory
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	// Create a bytes reader that implements io.ReaderAt
	bytesReader := bytes.NewReader(data)

	doc, err := document.Read(bytesReader, int64(len(data)))
	if err != nil {
		return "", err
	}
	defer doc.Close()

	var content strings.Builder
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			content.WriteString(run.Text())
		}
		content.WriteString("\n")
	}

	return content.String(), nil
}

func (s *DocumentService) extractTextContent(reader io.Reader) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

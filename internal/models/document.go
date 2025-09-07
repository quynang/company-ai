package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"not null;unique" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Document struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Content     string         `gorm:"type:text" json:"content"`
	Type        string         `gorm:"not null" json:"type"` // pdf, docx, txt
	Size        int64          `json:"size"`
	Categories  []Category     `gorm:"many2many:document_categories;" json:"categories,omitempty"`
	UploadedAt  time.Time      `json:"uploaded_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	ChunksCount int            `gorm:"-" json:"chunks_count"` // Virtual field for chunks count
}

type DocumentCategory struct {
	DocumentID uuid.UUID `gorm:"type:uuid;primary_key" json:"document_id"`
	CategoryID uuid.UUID `gorm:"type:uuid;primary_key" json:"category_id"`
	Document   Document  `gorm:"foreignKey:DocumentID" json:"document,omitempty"`
	Category   Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type DocumentChunk struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	DocumentID uuid.UUID `gorm:"type:uuid;not null" json:"document_id"`
	Document   Document  `gorm:"foreignKey:DocumentID" json:"document"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	ChunkIndex int       `gorm:"not null" json:"chunk_index"`
	Embedding  []float32 `gorm:"type:vector(768)" json:"-"` // Vector embedding for similarity search
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type User struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email      string         `gorm:"unique;not null" json:"email"`
	Name       string         `gorm:"not null" json:"name"`
	Department string         `json:"department"`
	Position   string         `json:"position"`
	StartDate  time.Time      `gorm:"not null" json:"start_date"`
	EmployeeID string         `gorm:"unique" json:"employee_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type ChatSession struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     *uuid.UUID `gorm:"type:uuid" json:"user_id"` // Optional user association
	User       *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CategoryID *uuid.UUID `gorm:"type:uuid" json:"category_id"` // Optional category filter
	Category   *Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name       string     `gorm:"not null" json:"name"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type ChatMessage struct {
	ID            uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SessionID     uuid.UUID   `gorm:"type:uuid;not null" json:"session_id"`
	Session       ChatSession `gorm:"foreignKey:SessionID" json:"session"`
	Role          string      `gorm:"not null" json:"role"` // user, assistant
	Content       string      `gorm:"type:text;not null" json:"content"`
	ContextChunks string      `gorm:"type:text" json:"-"` // Referenced document chunks as JSON string
	CreatedAt     time.Time   `json:"created_at"`
}

package models

import (
	"time"
	"github.com/google/uuid"
)

// ChatResponse represents the response structure for chat messages
type ChatResponse struct {
	Message    *ChatMessage `json:"message"`
	ActionCard *ActionCard  `json:"action_card,omitempty"`
}

// ActionCard represents an action that user can take
type ActionCard struct {
	Type        string      `json:"type"`        // "create_ticket", "contact_hr", etc.
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Action      ActionButton `json:"action"`
}

// ActionButton represents the button configuration
type ActionButton struct {
	Text     string            `json:"text"`
	Endpoint string            `json:"endpoint"`
	Method   string            `json:"method"`
	Payload  map[string]string `json:"payload,omitempty"`
}

// HRTicket represents an HR support ticket
type HRTicket struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	SessionID   uuid.UUID `gorm:"type:uuid;not null" json:"session_id"`
	Question    string    `gorm:"type:text;not null" json:"question"`
	Category    string    `gorm:"not null;default:'general'" json:"category"` // general, leave, policy, etc.
	Status      string    `gorm:"not null;default:'open'" json:"status"`      // open, in_progress, resolved, closed
	Priority    string    `gorm:"not null;default:'normal'" json:"priority"`  // low, normal, high, urgent
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

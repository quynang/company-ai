package services

import (
	"company-ai-training/internal/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketService struct {
	db *gorm.DB
}

func NewTicketService(db *gorm.DB) *TicketService {
	return &TicketService{
		db: db,
	}
}

// CreateTicket creates a new HR support ticket
func (s *TicketService) CreateTicket(sessionID uuid.UUID, userID *uuid.UUID, question, category, description string) (*models.HRTicket, error) {
	ticket := &models.HRTicket{
		ID:          uuid.New(),
		UserID:      userID,
		SessionID:   sessionID,
		Question:    question,
		Category:    category,
		Status:      "open",
		Priority:    "normal",
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(ticket).Error; err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	return ticket, nil
}

// GetTicket retrieves a ticket by ID
func (s *TicketService) GetTicket(ticketID uuid.UUID) (*models.HRTicket, error) {
	var ticket models.HRTicket
	if err := s.db.First(&ticket, "id = ?", ticketID).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

// GetUserTickets retrieves all tickets for a user
func (s *TicketService) GetUserTickets(userID uuid.UUID) ([]models.HRTicket, error) {
	var tickets []models.HRTicket
	if err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}

// GetSessionTickets retrieves all tickets for a chat session
func (s *TicketService) GetSessionTickets(sessionID uuid.UUID) ([]models.HRTicket, error) {
	var tickets []models.HRTicket
	if err := s.db.Where("session_id = ?", sessionID).Order("created_at DESC").Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}

// UpdateTicketStatus updates the status of a ticket
func (s *TicketService) UpdateTicketStatus(ticketID uuid.UUID, status string) error {
	return s.db.Model(&models.HRTicket{}).Where("id = ?", ticketID).Updates(map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}).Error
}

// GetAllTickets retrieves all tickets with pagination
func (s *TicketService) GetAllTickets(limit, offset int, status string) ([]models.HRTicket, error) {
	var tickets []models.HRTicket
	query := s.db.Order("created_at DESC")
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	if err := query.Limit(limit).Offset(offset).Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}


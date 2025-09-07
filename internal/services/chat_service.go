package services

import (
	"company-ai-training/internal/models"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatService struct {
	db            *gorm.DB
	vectorService *VectorService
	userService   *UserService
	geminiClient  *GeminiClientV2
}

func NewChatService(vectorService *VectorService, userService *UserService, geminiAPIKey string) *ChatService {
	return &ChatService{
		db:            vectorService.db,
		vectorService: vectorService,
		userService:   userService,
		geminiClient:  NewGeminiClientV2(geminiAPIKey),
	}
}

// CreateSession creates a new chat session
func (s *ChatService) CreateSession(name string, userID *uuid.UUID) (*models.ChatSession, error) {
	return s.CreateSessionWithCategory(name, userID, nil)
}

// CreateSessionWithCategory creates a new chat session with category filter
func (s *ChatService) CreateSessionWithCategory(name string, userID *uuid.UUID, categoryID *uuid.UUID) (*models.ChatSession, error) {
	session := &models.ChatSession{
		ID:         uuid.New(),
		UserID:     userID,
		CategoryID: categoryID,
		Name:       name,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.db.Create(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession retrieves a chat session
func (s *ChatService) GetSession(sessionID uuid.UUID) (*models.ChatSession, error) {
	var session models.ChatSession
	if err := s.db.Preload("Category").First(&session, "id = ?", sessionID).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// GetSessionMessages retrieves all messages for a session
func (s *ChatService) GetSessionMessages(sessionID uuid.UUID) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	if err := s.db.Where("session_id = ?", sessionID).Order("created_at ASC").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// GetAllSessions retrieves all chat sessions
func (s *ChatService) GetAllSessions() ([]models.ChatSession, error) {
	var sessions []models.ChatSession
	if err := s.db.Preload("Category").Order("updated_at DESC").Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// DeleteSession deletes a chat session and its messages
func (s *ChatService) DeleteSession(sessionID uuid.UUID) error {
	// Delete messages first
	if err := s.db.Where("session_id = ?", sessionID).Delete(&models.ChatMessage{}).Error; err != nil {
		return err
	}

	// Delete session
	return s.db.Delete(&models.ChatSession{}, "id = ?", sessionID).Error
}

// SendMessageWithResponse processes user message and generates AI response with action card support
func (s *ChatService) SendMessageWithResponse(sessionID uuid.UUID, userMessage string) (*models.ChatResponse, error) {
	// Save user message
	userMsg := &models.ChatMessage{
		ID:        uuid.New(),
		SessionID: sessionID,
		Role:      "user",
		Content:   userMessage,
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(userMsg).Error; err != nil {
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}

	// Get session to check for category filter
	session, err := s.GetSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Search for relevant document chunks with category filter
	relevantChunks, err := s.vectorService.SearchSimilarChunksWithCategory(userMessage, 5, session.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to search relevant chunks: %w", err)
	}

	// Build context from relevant chunks
	var contextParts []string
	var contextChunkIDs []uuid.UUID
	hasRelevantInfo := false

	for _, chunk := range relevantChunks {
		contextParts = append(contextParts, fmt.Sprintf("Document: %s\nContent: %s", chunk.Document.Name, chunk.Content))
		contextChunkIDs = append(contextChunkIDs, chunk.ID)
		hasRelevantInfo = true
	}

	context := strings.Join(contextParts, "\n\n---\n\n")

	// Get conversation history
	messages, err := s.GetSessionMessages(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation history: %w", err)
	}

	// Build conversation for AI
	var conversation []Message

	// Get user context if session has user
	var userContext string
	if session.UserID != nil {
		if user, err := s.userService.GetUser(*session.UserID); err == nil {
			userContext, _ = s.userService.GetUserContext(user)
		}
	}

	// Add system prompt with context
	systemPrompt := s.buildSystemPrompt(context, userContext)
	conversation = append(conversation, Message{
		Role:    "system",
		Content: systemPrompt,
	})

	// Add recent conversation history (last 10 messages)
	historyStart := len(messages) - 10
	if historyStart < 0 {
		historyStart = 0
	}

	for _, msg := range messages[historyStart:] {
		if msg.Role == "user" || msg.Role == "assistant" {
			conversation = append(conversation, Message{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	}

	// Build conversation string for Gemini
	conversationText := ""
	for _, msg := range conversation {
		conversationText += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content)
	}

	// Generate AI response
	response, err := s.geminiClient.Chat(conversationText)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI response: %w", err)
	}

	// Convert context chunks to JSON string
	contextChunksJSON, _ := json.Marshal(contextChunkIDs)

	// Save assistant message
	assistantMsg := &models.ChatMessage{
		ID:            uuid.New(),
		SessionID:     sessionID,
		Role:          "assistant",
		Content:       response,
		ContextChunks: string(contextChunksJSON),
		CreatedAt:     time.Now(),
	}

	if err := s.db.Create(assistantMsg).Error; err != nil {
		return nil, fmt.Errorf("failed to save assistant message: %w", err)
	}

	// Update session timestamp
	s.db.Model(&models.ChatSession{}).Where("id = ?", sessionID).Update("updated_at", time.Now())

	// Create response
	chatResponse := &models.ChatResponse{
		Message: assistantMsg,
	}

	// If no relevant information found, add action card
	if !hasRelevantInfo || strings.Contains(strings.ToLower(response), "không tìm thấy") ||
		strings.Contains(strings.ToLower(response), "không có thông tin") {
		chatResponse.ActionCard = &models.ActionCard{
			Type:        "create_ticket",
			Title:       "Không tìm thấy thông tin?",
			Description: "Tạo ticket để được hỗ trợ trực tiếp từ bộ phận nhân sự",
			Action: models.ActionButton{
				Text:     "Tạo Ticket Hỏi HR",
				Endpoint: "/api/v1/tickets",
				Method:   "POST",
				Payload: map[string]string{
					"question":   userMessage,
					"category":   "general",
					"session_id": sessionID.String(),
				},
			},
		}
	}

	return chatResponse, nil
}

// SendMessage processes a user message and generates AI response (backward compatibility)
func (s *ChatService) SendMessage(sessionID uuid.UUID, userMessage string) (*models.ChatMessage, error) {
	response, err := s.SendMessageWithResponse(sessionID, userMessage)
	if err != nil {
		return nil, err
	}
	return response.Message, nil
}

// buildSystemPrompt creates system prompt with document context and user info
func (s *ChatService) buildSystemPrompt(context string, userContext string) string {
	prompt := `Bạn là một AI assistant được thiết kế để trả lời câu hỏi dựa trên tài liệu nội bộ của công ty.

## HƯỚNG DẪN:
1. Chỉ trả lời dựa trên thông tin có trong tài liệu được cung cấp.
2. Sử dụng thông tin cá nhân của nhân viên để tính toán và trả lời chính xác.
3. Nếu không tìm thấy thông tin trong tài liệu, hãy trả lời: 
   **"Tôi không tìm thấy thông tin này trong tài liệu."**
4. Trả lời bằng tiếng Việt một cách chính xác, dễ hiểu và có cấu trúc rõ ràng.
5. Khi tính toán ngày phép, hãy sử dụng thông tin thâm niên của nhân viên.
6. Nếu có nhiều nguồn thông tin, hãy tổng hợp và trình bày một cách logic.
7. **Luôn định dạng câu trả lời bằng Markdown**:
   - Sử dụng heading (##) cho các phần chính.
   - Dùng danh sách đánh số hoặc bullet (-) cho từng ý.
   - Nếu có link tài liệu đính kèm, trình bày theo dạng: Link tải [Tên hiển thị](URL).
   - Nếu có bảng, hãy sử dụng bảng Markdown chuẩn.
`

	if userContext != "" {
		prompt += fmt.Sprintf(`

## NGỮ CẢNH NGƯỜI DÙNG:
%s`, userContext)
	}

	if context != "" {
		prompt += fmt.Sprintf(`

## TÀI LIỆU THAM KHẢO:
%s`, context)
	}

	prompt += `

---

## YÊU CẦU:
Hãy trả lời câu hỏi dựa trên thông tin trên và xuất kết quả ở định dạng Markdown.`

	fmt.Println(prompt)

	return prompt
}

package api

import (
	"company-ai-training/internal/services"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handlers struct {
	documentService *services.DocumentService
	vectorService   *services.VectorService
	chatService     *services.ChatService
	userService     *services.UserService
	ticketService   *services.TicketService
}

func NewHandlers(docService *services.DocumentService, vecService *services.VectorService, chatService *services.ChatService, userService *services.UserService, ticketService *services.TicketService) *Handlers {
	return &Handlers{
		documentService: docService,
		vectorService:   vecService,
		chatService:     chatService,
		userService:     userService,
		ticketService:   ticketService,
	}
}

// Document handlers

func (h *Handlers) UploadDocument(c *gin.Context) {
	// Check if it's a file upload or text content
	file, err := c.FormFile("file")
	if err != nil {
		// Try to get text content from different sources
		var content, name string

		// First try form data
		content = c.PostForm("content")
		name = c.PostForm("name")

		// If not found, try JSON body
		if content == "" || name == "" {
			var jsonData struct {
				Name    string `json:"name"`
				Content string `json:"content"`
			}

			if err := c.ShouldBindJSON(&jsonData); err == nil {
				content = jsonData.Content
				name = jsonData.Name
			}
		}

		if content == "" || name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Either file or content+name must be provided"})
			return
		}

		// Create document from text content
		doc, err := h.documentService.CreateDocumentFromText(name, content)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Process document for vector search using semantic chunking
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Panic in embedding process for document %s: %v\n", doc.Name, r)
				}
			}()

			// Use semantic chunking by default
			config := services.DefaultChunkConfig()
			if err := h.vectorService.ChunkAndEmbedDocumentWithSemantics(doc, config); err != nil {
				// Fallback to legacy chunking if semantic fails
				fmt.Printf("Semantic chunking failed for document %s, falling back to legacy: %v\n", doc.Name, err)
				if err := h.vectorService.ChunkAndEmbedDocument(doc); err != nil {
					fmt.Printf("Error embedding document %s: %v\n", doc.Name, err)
				} else {
					fmt.Printf("Successfully embedded document %s with legacy chunking\n", doc.Name)
				}
			} else {
				fmt.Printf("Successfully embedded document %s with semantic chunking\n", doc.Name)
			}
		}()

		c.JSON(http.StatusCreated, gin.H{
			"message":  "Document created successfully",
			"document": doc,
		})
		return
	}

	// Upload document from file
	doc, err := h.documentService.UploadDocument(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Process document for vector search using semantic chunking
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Panic in embedding process for document %s: %v\n", doc.Name, r)
			}
		}()

		// Use semantic chunking by default
		config := services.DefaultChunkConfig()
		if err := h.vectorService.ChunkAndEmbedDocumentWithSemantics(doc, config); err != nil {
			// Fallback to legacy chunking if semantic fails
			fmt.Printf("Semantic chunking failed for document %s, falling back to legacy: %v\n", doc.Name, err)
			if err := h.vectorService.ChunkAndEmbedDocument(doc); err != nil {
				fmt.Printf("Error embedding document %s: %v\n", doc.Name, err)
			} else {
				fmt.Printf("Successfully embedded document %s with legacy chunking\n", doc.Name)
			}
		} else {
			fmt.Printf("Successfully embedded document %s with semantic chunking\n", doc.Name)
		}
	}()

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Document uploaded successfully",
		"document": doc,
	})
}

func (h *Handlers) GetDocuments(c *gin.Context) {
	docs, err := h.documentService.GetAllDocuments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"documents": docs})
}

func (h *Handlers) GetDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	doc, err := h.documentService.GetDocument(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"document": doc})
}

func (h *Handlers) DeleteDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	if err := h.documentService.DeleteDocument(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}

// Search handlers

func (h *Handlers) SearchDocuments(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	chunks, err := h.vectorService.SearchSimilarChunks(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":   query,
		"results": chunks,
	})
}

// Chat handlers

// User handlers

type CreateUserRequest struct {
	Email      string `json:"email" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Department string `json:"department"`
	Position   string `json:"position"`
	EmployeeID string `json:"employee_id"`
	StartDate  string `json:"start_date" binding:"required"` // Format: "2006-01-02"
}

func (h *Handlers) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse start date
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
		return
	}

	user, err := h.userService.CreateUser(req.Email, req.Name, req.Department, req.Position, req.EmployeeID, startDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func (h *Handlers) GetUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *Handlers) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Calculate vacation days

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *Handlers) GetUserByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email parameter is required"})
		return
	}

	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Calculate vacation days

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

type CreateSessionRequest struct {
	Name   string     `json:"name" binding:"required"`
	UserID *uuid.UUID `json:"user_id,omitempty"` // Optional user association
}

func (h *Handlers) CreateChatSession(c *gin.Context) {
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.chatService.CreateSession(req.Name, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"session": session})
}

func (h *Handlers) GetChatSessions(c *gin.Context) {
	sessions, err := h.chatService.GetAllSessions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

func (h *Handlers) GetChatSession(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	session, err := h.chatService.GetSession(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	messages, err := h.chatService.GetSessionMessages(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session":  session,
		"messages": messages,
	})
}

func (h *Handlers) DeleteChatSession(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	if err := h.chatService.DeleteSession(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session deleted successfully"})
}

type SendMessageRequest struct {
	Message string `json:"message" binding:"required"`
}

func (h *Handlers) SendMessage(c *gin.Context) {
	idStr := c.Param("id")
	sessionID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.chatService.SendMessageWithResponse(sessionID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}

// Ticket handlers

type CreateTicketRequest struct {
	Question    string `json:"question" binding:"required"`
	Category    string `json:"category"`
	Description string `json:"description"`
	SessionID   string `json:"session_id" binding:"required"`
}

func (h *Handlers) CreateTicket(c *gin.Context) {
	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	// Default category if not provided
	if req.Category == "" {
		req.Category = "general"
	}

	ticket, err := h.ticketService.CreateTicket(sessionID, nil, req.Question, req.Category, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"ticket":  ticket,
		"message": "Ticket đã được tạo thành công! Bộ phận nhân sự sẽ liên hệ với bạn sớm.",
	})
}

func (h *Handlers) GetTickets(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	status := c.Query("status")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	tickets, err := h.ticketService.GetAllTickets(limit, offset, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tickets": tickets})
}

func (h *Handlers) GetTicket(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	ticket, err := h.ticketService.GetTicket(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ticket": ticket})
}

type UpdateTicketStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (h *Handlers) UpdateTicketStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var req UpdateTicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.ticketService.UpdateTicketStatus(id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket status updated successfully"})
}

// Admin Dashboard handlers

type UpdateDocumentRequest struct {
	Content string `json:"content" binding:"required"`
}

func (h *Handlers) UpdateDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var req UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update document content
	if err := h.documentService.UpdateDocument(id, req.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get updated document
	doc, err := h.documentService.GetDocument(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated document"})
		return
	}

	// Re-embed the updated document
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Panic in re-embedding process for document %s: %v\n", doc.Name, r)
			}
		}()

		// Delete old chunks first
		if err := h.vectorService.DeleteDocumentChunks(doc.ID); err != nil {
			fmt.Printf("Error deleting old chunks for document %s: %v\n", doc.Name, err)
		} else {
			fmt.Printf("Deleted old chunks for document %s\n", doc.Name)
		}

		// Create new chunks and embeddings using semantic chunking
		config := services.DefaultChunkConfig()
		if err := h.vectorService.ChunkAndEmbedDocumentWithSemantics(doc, config); err != nil {
			// Fallback to legacy chunking if semantic fails
			fmt.Printf("Semantic chunking failed for document %s, falling back to legacy: %v\n", doc.Name, err)
			if err := h.vectorService.ChunkAndEmbedDocument(doc); err != nil {
				fmt.Printf("Error re-embedding document %s: %v\n", doc.Name, err)
			} else {
				fmt.Printf("Successfully re-embedded document %s with legacy chunking\n", doc.Name)
			}
		} else {
			fmt.Printf("Successfully re-embedded document %s with semantic chunking\n", doc.Name)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message":  "Document updated and re-embedding started",
		"document": doc,
	})
}

func (h *Handlers) ReembedDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	// Get document
	doc, err := h.documentService.GetDocument(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// Re-embed the document
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Panic in re-embedding process for document %s: %v\n", doc.Name, r)
			}
		}()

		// Delete old chunks first
		if err := h.vectorService.DeleteDocumentChunks(doc.ID); err != nil {
			fmt.Printf("Error deleting old chunks for document %s: %v\n", doc.Name, err)
		} else {
			fmt.Printf("Deleted old chunks for document %s\n", doc.Name)
		}

		// Create new chunks and embeddings using semantic chunking
		config := services.DefaultChunkConfig()
		if err := h.vectorService.ChunkAndEmbedDocumentWithSemantics(doc, config); err != nil {
			// Fallback to legacy chunking if semantic fails
			fmt.Printf("Semantic chunking failed for document %s, falling back to legacy: %v\n", doc.Name, err)
			if err := h.vectorService.ChunkAndEmbedDocument(doc); err != nil {
				fmt.Printf("Error re-embedding document %s: %v\n", doc.Name, err)
			} else {
				fmt.Printf("Successfully re-embedded document %s with legacy chunking\n", doc.Name)
			}
		} else {
			fmt.Printf("Successfully re-embedded document %s with semantic chunking\n", doc.Name)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message":  "Document re-embedding started",
		"document": doc,
	})
}

// Semantic chunking endpoint
type SemanticChunkingRequest struct {
	DocumentID string `json:"document_id" binding:"required"`
	Config     *services.ChunkConfig `json:"config,omitempty"`
}

func (h *Handlers) ReembedWithSemanticChunking(c *gin.Context) {
	var req SemanticChunkingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse document ID
	docID, err := uuid.Parse(req.DocumentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	// Get document
	doc, err := h.documentService.GetDocument(docID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// Use provided config or default
	config := req.Config
	if config == nil {
		config = services.DefaultChunkConfig()
	}

	// Re-embed with semantic chunking
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Panic in semantic re-embedding process for document %s: %v\n", doc.Name, r)
			}
		}()

		// Delete old chunks first
		if err := h.vectorService.DeleteDocumentChunks(doc.ID); err != nil {
			fmt.Printf("Error deleting old chunks for document %s: %v\n", doc.Name, err)
		} else {
			fmt.Printf("Deleted old chunks for document %s\n", doc.Name)
		}

		// Create new chunks and embeddings using semantic chunking
		if err := h.vectorService.ChunkAndEmbedDocumentWithSemantics(doc, config); err != nil {
			fmt.Printf("Error semantic re-embedding document %s: %v\n", doc.Name, err)
		} else {
			fmt.Printf("Successfully semantic re-embedded document %s\n", doc.Name)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message":  "Document semantic re-embedding started",
		"document": doc,
		"config":   config,
	})
}

// Health check

func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "company-ai-training",
	})
}

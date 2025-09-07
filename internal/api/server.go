package api

import (
	"company-ai-training/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router   *gin.Engine
	handlers *Handlers
}

func NewServer(docService *services.DocumentService, vecService *services.VectorService, chatService *services.ChatService, userService *services.UserService, ticketService *services.TicketService) *Server {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// CORS middleware - must be first, before any other middleware
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://0.0.0.0:3000", "http://10.67.21.180:3000", "http://192.168.1.100:3000", "*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
		AllowCredentials: false,
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Methods", "Access-Control-Allow-Headers"},
		MaxAge:           86400,
	}
	router.Use(cors.New(config))

	// Add logger and recovery middleware after CORS
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Initialize handlers
	handlers := NewHandlers(docService, vecService, chatService, userService, ticketService)

	server := &Server{
		router:   router,
		handlers: handlers,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api/v1")

	// Health check
	api.GET("/health", s.handlers.HealthCheck)

	// Document routes
	documents := api.Group("/documents")
	{
		documents.POST("/upload", s.handlers.UploadDocument)
		documents.GET("/", s.handlers.GetDocuments)
		documents.GET("", s.handlers.GetDocuments) // Add route without trailing slash
		documents.GET("/:id", s.handlers.GetDocument)
		documents.PUT("/:id", s.handlers.UpdateDocument)
		documents.DELETE("/:id", s.handlers.DeleteDocument)
		documents.POST("/:id/reembed", s.handlers.ReembedDocument)
		documents.POST("/semantic-reembed", s.handlers.ReembedWithSemanticChunking)
	}

	// Search routes
	search := api.Group("/search")
	{
		search.GET("/", s.handlers.SearchDocuments)
	}

	// User routes
	users := api.Group("/users")
	{
		users.POST("/", s.handlers.CreateUser)
		users.GET("/", s.handlers.GetUsers)
		users.GET("/:id", s.handlers.GetUser)
		users.GET("/by-email", s.handlers.GetUserByEmail)
	}

	// Chat routes
	chat := api.Group("/chat")
	{
		chat.POST("/sessions", s.handlers.CreateChatSession)
		chat.GET("/sessions", s.handlers.GetChatSessions)
		chat.GET("/sessions/:id", s.handlers.GetChatSession)
		chat.DELETE("/sessions/:id", s.handlers.DeleteChatSession)
		chat.POST("/sessions/:id/messages", s.handlers.SendMessage)
	}

	// Ticket routes
	tickets := api.Group("/tickets")
	{
		tickets.POST("/", s.handlers.CreateTicket)
		tickets.GET("/", s.handlers.GetTickets)
		tickets.GET("/:id", s.handlers.GetTicket)
		tickets.PUT("/:id/status", s.handlers.UpdateTicketStatus)
	}
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

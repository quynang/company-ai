package main

import (
	"company-ai-training/internal/api"
	"company-ai-training/internal/config"
	"company-ai-training/internal/database"
	"company-ai-training/internal/services"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize services
	documentService := services.NewDocumentService(db)
	vectorService := services.NewVectorService(db, cfg.GeminiAPIKey)
	userService := services.NewUserService(db)
	chatService := services.NewChatService(vectorService, userService, cfg.GeminiAPIKey)
	ticketService := services.NewTicketService(db)
	categoryService := services.NewCategoryService(db)

	// Initialize API server
	server := api.NewServer(documentService, vectorService, chatService, userService, ticketService, categoryService)

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := server.Start(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

# Company AI Training - Makefile
# ================================

# Include .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

.PHONY: help run start-db stop-db restart-db logs status env

# Default target
help: ## Show this help message
	@echo "Company AI Training - Available Commands:"
	@echo "=========================================="
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Main Commands
run: ## Run the backend server
	@echo "🚀 Starting Company AI Training backend..."
	@echo "📋 Make sure GEMINI_API_KEY is set and database is running"
	@echo "🌐 Backend will be available at: http://localhost:8082"
	@echo ""
	PORT=8082 go run main.go

# Database Management
start-db: ## Start PostgreSQL database
	@echo "🐘 Starting PostgreSQL database..."
	docker-compose up -d postgres
	@echo "✅ Database started on port 5434"

stop-db: ## Stop PostgreSQL database
	@echo "🛑 Stopping PostgreSQL database..."
	docker-compose down
	@echo "✅ Database stopped"

restart-db: ## Restart PostgreSQL database
	@echo "🔄 Restarting PostgreSQL database..."
	docker-compose restart postgres
	@echo "✅ Database restarted"

# Utility Commands
env: ## Create .env file from template
	@echo "📝 Creating .env file from template..."
	@if [ ! -f .env ]; then \
		cp env.example .env; \
		echo "✅ .env file created from env.example"; \
		echo "📝 Please edit .env file with your GEMINI_API_KEY"; \
	else \
		echo "⚠️  .env file already exists"; \
	fi

logs: ## Show database logs
	@echo "📋 Showing database logs..."
	docker-compose logs -f postgres

status: ## Show application status
	@echo "📊 Application Status:"
	@echo "====================="
	@echo "Database:"
	@docker-compose ps postgres
	@echo ""
	@echo "Environment Variables:"
	@echo "GEMINI_API_KEY: $(if $(GEMINI_API_KEY),✅ Set,❌ Not set)"
	@echo "DATABASE_URL: $(if $(DATABASE_URL),✅ Set,❌ Not set)"

# Default environment variables
#!/bin/bash

echo "🚀 Company AI Training System Setup"
echo "=================================="

# Check if docker and docker-compose are installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Create .env file if not exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp env.example .env
    echo "✅ .env file created. Please review and modify if needed."
fi

# Start the services
echo "🐳 Starting Docker services..."
docker-compose up -d

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 10

# Check if Ollama is ready
echo "🤖 Checking Ollama service..."
max_attempts=30
attempt=1

while [ $attempt -le $max_attempts ]; do
    if curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
        echo "✅ Ollama is ready!"
        break
    else
        echo "Waiting for Ollama... (attempt $attempt/$max_attempts)"
        sleep 5
        ((attempt++))
    fi
done

if [ $attempt -gt $max_attempts ]; then
    echo "❌ Ollama failed to start. Please check logs: docker-compose logs ollama"
    exit 1
fi

# Pull required models
echo "📥 Downloading AI models (this may take a while)..."

echo "Downloading chat model (llama2)..."
docker exec company_ai_ollama ollama pull llama2

echo "Downloading embedding model (nomic-embed-text)..."
docker exec company_ai_ollama ollama pull nomic-embed-text

# Check if app is running
echo "🔍 Checking application health..."
sleep 5

if curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo "✅ Application is running successfully!"
else
    echo "⚠️  Application might still be starting. Check logs: docker-compose logs app"
fi

echo ""
echo "🎉 Setup completed!"
echo ""
echo "📋 Next steps:"
echo "1. Check application health: curl http://localhost:8080/api/v1/health"
echo "2. Upload a document: curl -X POST -F 'file=@document.pdf' http://localhost:8080/api/v1/documents/upload"
echo "3. Create a chat session: curl -X POST -H 'Content-Type: application/json' -d '{\"name\":\"Test Chat\"}' http://localhost:8080/api/v1/chat/sessions"
echo ""
echo "📚 For more information, see README.md"
echo "🔧 To stop services: docker-compose down"
echo "📊 To view logs: docker-compose logs -f"

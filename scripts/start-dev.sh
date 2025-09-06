#!/bin/bash

# Script để khởi động cả backend và frontend cho development

echo "🚀 Starting Company AI Chat Development Environment"
echo "================================================="

# Kiểm tra Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Kiểm tra Node.js
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js first."
    exit 1
fi

# Chuyển đến thư mục root
cd "$(dirname "$0")/.."

echo "🐳 Starting backend services (PostgreSQL, Ollama, Go API)..."
echo "This may take a few minutes on first run..."

# Khởi động backend services
docker-compose up -d

# Chờ services khởi động
echo "⏳ Waiting for services to be ready..."
sleep 10

# Kiểm tra backend health
echo "🔍 Checking backend health..."
for i in {1..30}; do
    if curl -s http://localhost:8082/api/v1/health > /dev/null 2>&1; then
        echo "✅ Backend is ready!"
        break
    fi
    
    if [ $i -eq 30 ]; then
        echo "❌ Backend failed to start. Please check Docker logs:"
        echo "   docker-compose logs app"
        exit 1
    fi
    
    echo "   Attempt $i/30 - Backend not ready yet..."
    sleep 2
done

# Chuyển đến thư mục frontend
cd frontend

# Kiểm tra và cài đặt dependencies
if [ ! -d "node_modules" ]; then
    echo "📦 Installing frontend dependencies..."
    npm install
fi

# Tạo .env file nếu chưa có
if [ ! -f ".env" ]; then
    echo "⚙️ Creating .env file..."
    cp env.example .env
fi

echo ""
echo "🎉 Development environment is ready!"
echo "================================================="
echo "Backend API: http://localhost:8082/api/v1"
echo "Frontend:    http://localhost:3000 (will open automatically)"
echo ""
echo "Services running:"
echo "- PostgreSQL: localhost:5434"
echo "- Ollama:     localhost:11435"
echo "- Backend:    localhost:8082"
echo ""
echo "To stop all services, run:"
echo "  docker-compose down"
echo ""
echo "Starting frontend development server..."
echo "Press Ctrl+C to stop the frontend server"
echo ""

# Khởi động frontend
npm start



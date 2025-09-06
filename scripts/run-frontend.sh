#!/bin/bash

# Script để chạy React frontend

echo "🚀 Starting Company AI Chat Frontend..."

# Chuyển đến thư mục frontend
cd "$(dirname "$0")/../frontend"

# Kiểm tra xem đã cài đặt dependencies chưa
if [ ! -d "node_modules" ]; then
    echo "📦 Installing dependencies..."
    npm install
fi

# Kiểm tra file .env
if [ ! -f ".env" ]; then
    if [ -f "env.example" ]; then
        echo "⚙️ Creating .env file from env.example..."
        cp env.example .env
        echo "✅ .env file created. Please review and update if needed."
    else
        echo "⚠️ Warning: No .env file found. Creating default one..."
        echo "REACT_APP_API_URL=http://localhost:8080/api/v1" > .env
    fi
fi

# Chạy ứng dụng
echo "🎉 Starting development server..."
echo "Frontend will be available at: http://localhost:3000"
echo "Make sure backend is running at: http://localhost:8080"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

npm start

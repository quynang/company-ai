#!/bin/bash

echo "🚀 Starting Company AI Training System for Public Access..."
echo ""

# Get local IP address
LOCAL_IP=$(hostname -I | awk '{print $1}')
echo "📍 Your local IP address: $LOCAL_IP"
echo "🌐 Frontend will be available at: http://$LOC_IP:3000"
echo "🔧 Backend API will be available at: http://$LOC_IP:8082"
echo ""

# Start backend
echo "🔧 Starting backend services..."
docker-compose up -d

# Wait for backend to be ready
echo "⏳ Waiting for backend to be ready..."
sleep 10

# Check if backend is running
if curl -s http://localhost:8082/api/v1/health > /dev/null 2>&1; then
    echo "✅ Backend is running successfully"
else
    echo "❌ Backend failed to start. Check logs with: docker-compose logs app"
    exit 1
fi

echo ""
echo "🎯 Starting frontend with public access..."
echo "📱 Other users in the same WiFi can access at: http://$LOCAL_IP:3000"
echo ""

# Set environment variable for frontend
export HOST=0.0.0.0
export REACT_APP_API_URL=http://$LOCAL_IP:8082/api/v1

# Start frontend
cd frontend && npm start

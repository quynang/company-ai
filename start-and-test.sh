#!/bin/bash

echo "🚀 Starting Company AI Training System for Public Access..."
echo ""

# Get local IP address
LOCAL_IP=$(hostname -I | awk '{print $1}')
echo "📍 Your local IP address: $LOCAL_IP"
echo "🌐 Frontend will be available at: http://$LOCAL_IP:3000"
echo "🔧 Backend API will be available at: http://$LOCAL_IP:8082"
echo ""

# Start backend
echo "🔧 Starting backend services..."
docker-compose down
docker-compose up --build -d

# Wait for backend to be ready
echo "⏳ Waiting for backend to be ready..."
sleep 15

# Check if backend is running
echo "🔍 Testing backend connectivity..."
if curl -s http://localhost:8082/api/v1/health > /dev/null 2>&1; then
    echo "✅ Backend is running successfully on localhost"
else
    echo "❌ Backend failed to start on localhost"
    docker-compose logs app --tail=10
    exit 1
fi

# Test backend from local IP
if curl -s http://$LOCAL_IP:8082/api/v1/health > /dev/null 2>&1; then
    echo "✅ Backend is accessible from local IP: $LOCAL_IP"
else
    echo "❌ Backend is NOT accessible from local IP: $LOCAL_IP"
    echo "   This means other devices cannot call the API"
fi

echo ""
echo "🎯 Starting frontend with public access..."
echo "📱 Other users in the same WiFi can access at: http://$LOCAL_IP:3000"
echo ""

# Set environment variable for frontend
export HOST=0.0.0.0
export REACT_APP_API_URL=http://$LOCAL_IP:8082/api/v1

echo "🔧 Frontend will use API at: $REACT_APP_API_URL"
echo ""

# Start frontend
cd frontend && npm start

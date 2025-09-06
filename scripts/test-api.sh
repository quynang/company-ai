#!/bin/bash

# Script để test API connection

echo "🔍 Testing API Connection..."
echo "============================"

API_URL="http://localhost:8082/api/v1"

# Test health endpoint
echo "1. Testing health endpoint..."
if curl -s "$API_URL/health" | grep -q "ok"; then
    echo "✅ Health check passed"
else
    echo "❌ Health check failed"
    echo "Make sure backend is running: docker-compose up -d"
    exit 1
fi

# Test chat sessions endpoint
echo "2. Testing chat sessions endpoint..."
if curl -s -X GET "$API_URL/chat/sessions" > /dev/null; then
    echo "✅ Chat sessions endpoint accessible"
else
    echo "❌ Chat sessions endpoint failed"
    exit 1
fi

# Test CORS headers
echo "3. Testing CORS headers..."
CORS_HEADERS=$(curl -s -I -X OPTIONS "$API_URL/chat/sessions" | grep -i "access-control")
if [ ! -z "$CORS_HEADERS" ]; then
    echo "✅ CORS headers present"
    echo "   $CORS_HEADERS"
else
    echo "⚠️ CORS headers not found"
fi

echo ""
echo "🎉 API is ready for frontend connection!"
echo "Frontend should connect to: $API_URL"



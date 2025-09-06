#!/bin/bash

echo "🚀 Testing Google Gemini API Integration"
echo "========================================"

# Check if GEMINI_API_KEY is set
if [ -z "$GEMINI_API_KEY" ]; then
    echo "❌ GEMINI_API_KEY environment variable is not set"
    echo "Please set your Gemini API key:"
    echo "export GEMINI_API_KEY=your_api_key_here"
    exit 1
fi

echo "✅ GEMINI_API_KEY is set"

# Test Gemini Chat API
echo ""
echo "🧪 Testing Gemini Chat API..."
CHAT_RESPONSE=$(curl -s -H "Content-Type: application/json" \
     -d '{"contents":[{"parts":[{"text":"Hello! Can you respond with just the word SUCCESS if you can understand me?"}]}]}' \
     "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=$GEMINI_API_KEY")

if echo "$CHAT_RESPONSE" | grep -q "SUCCESS"; then
    echo "✅ Gemini Chat API is working!"
else
    echo "❌ Gemini Chat API test failed"
    echo "Response: $CHAT_RESPONSE"
    exit 1
fi

# Test Gemini Embedding API
echo ""
echo "🧪 Testing Gemini Embedding API..."
EMBED_RESPONSE=$(curl -s -H "Content-Type: application/json" \
     -d '{"model":"models/embedding-001","content":{"parts":[{"text":"Hello world"}]},"taskType":"RETRIEVAL_DOCUMENT"}' \
     "https://generativelanguage.googleapis.com/v1beta/models/embedding-001:embedContent?key=$GEMINI_API_KEY")

if echo "$EMBED_RESPONSE" | grep -q "embedding"; then
    echo "✅ Gemini Embedding API is working!"
else
    echo "❌ Gemini Embedding API test failed"
    echo "Response: $EMBED_RESPONSE"
    exit 1
fi

echo ""
echo "🎉 All Gemini API tests passed!"
echo "Your system is ready to use Google Gemini API."
echo ""
echo "Next steps:"
echo "1. Make sure PostgreSQL is running: docker-compose up -d postgres"
echo "2. Start the application: go run main.go"
echo "3. Test the full system with: ./scripts/test-api.sh"

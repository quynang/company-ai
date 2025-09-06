#!/bin/bash

echo "🔍 Testing connection to Company AI Training System..."
echo ""

# Get local IP address
LOCAL_IP=$(hostname -I | awk '{print $1}')
echo "📍 Your local IP address: $LOCAL_IP"
echo ""

# Test backend connection
echo "🔧 Testing backend connection..."
if curl -s http://$LOCAL_IP:8082/api/v1/health > /dev/null 2>&1; then
    echo "✅ Backend is accessible at: http://$LOCAL_IP:8082"
else
    echo "❌ Backend is NOT accessible at: http://$LOCAL_IP:8082"
    echo "   Trying localhost..."
    if curl -s http://localhost:8082/api/v1/health > /dev/null 2>&1; then
        echo "✅ Backend is accessible at localhost:8082"
    else
        echo "❌ Backend is not running"
    fi
fi

echo ""

# Test frontend connection
echo "🌐 Testing frontend connection..."
if curl -s http://$LOCAL_IP:3000 > /dev/null 2>&1; then
    echo "✅ Frontend is accessible at: http://$LOCAL_IP:3000"
else
    echo "❌ Frontend is NOT accessible at: http://$LOCAL_IP:3000"
    echo "   Trying localhost..."
    if curl -s http://localhost:3000 > /dev/null 2>&1; then
        echo "✅ Frontend is accessible at localhost:3000"
    else
        echo "❌ Frontend is not running"
    fi
fi

echo ""

# Check if processes are running
echo "📊 Checking running processes..."
echo "Backend (port 8082):"
netstat -tlnp | grep :8082 || echo "   Not found"

echo "Frontend (port 3000):"
netstat -tlnp | grep :3000 || echo "   Not found"

echo ""
echo "🎯 URLs to share with others:"
echo "   Frontend: http://$LOCAL_IP:3000"
echo "   Backend:  http://$LOCAL_IP:8082"
echo ""
echo "💡 If others can't access, check:"
echo "   1. Same WiFi network"
echo "   2. Firewall settings"
echo "   3. Router settings"
echo "   4. Try from another device in the same network"

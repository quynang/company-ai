#!/bin/bash

# Get local IP address
LOCAL_IP=$(hostname -I | awk '{print $1}')

echo "Starting frontend with public IP: $LOCAL_IP"
echo "Other users can access at: http://$LOCAL_IP:3000"
echo ""

# Set environment variable to bind to all interfaces
export HOST=0.0.0.0

# Start React app
npm start

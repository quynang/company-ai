#!/bin/bash

# Script upload từng chunk một cách tuần tự để tránh overload

echo "🚀 Uploading document chunks..."
echo "================================="

CHUNK_DIR="chunked_documents"
API_URL="http://localhost:8082/api/v1/documents/upload"

if [ ! -d "$CHUNK_DIR" ]; then
    echo "❌ Error: Directory $CHUNK_DIR not found"
    echo "Please run: python3 scripts/chunk-document.py <input_file> first"
    exit 1
fi

# Đếm số file chunks
total_files=$(find "$CHUNK_DIR" -name "*.txt" | wc -l)
echo "📄 Found $total_files chunk files"
echo ""

counter=0
success=0
failed=0

# Upload từng file với delay giữa các request
for file in "$CHUNK_DIR"/*.txt; do
    if [ -f "$file" ]; then
        counter=$((counter + 1))
        filename=$(basename "$file")
        
        echo "[$counter/$total_files] Uploading: $filename"
        
        # Upload file
        response=$(curl -s -X POST -F "file=@$file" "$API_URL")
        
        # Kiểm tra kết quả
        if echo "$response" | grep -q "Document uploaded successfully"; then
            echo "✅ Success"
            success=$((success + 1))
        else
            echo "❌ Failed: $response"
            failed=$((failed + 1))
        fi
        
        # Delay 2 giây giữa các upload để tránh overload
        if [ $counter -lt $total_files ]; then
            echo "   Waiting 2 seconds..."
            sleep 2
        fi
        
        echo ""
    fi
done

echo "================================="
echo "📊 Upload Summary:"
echo "   Total files: $total_files"
echo "   Successful: $success"
echo "   Failed: $failed"

if [ $failed -eq 0 ]; then
    echo "🎉 All chunks uploaded successfully!"
    echo ""
    echo "Testing search..."
    sleep 5
    curl -s "http://localhost:8082/api/v1/search/?q=phep%20nam&limit=2" | grep -o '"results":\[.*\]' | head -1
else
    echo "⚠️  Some uploads failed. Please check the errors above."
fi


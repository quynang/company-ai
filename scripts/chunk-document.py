#!/usr/bin/env python3
"""
Script để chia nhỏ tài liệu thành các chunks nhỏ hơn
để giảm tải cho quá trình embedding
"""

import os
import sys
import re

def chunk_text(text, chunk_size=200, overlap=50):
    """Chia văn bản thành các chunks nhỏ với overlap"""
    words = text.split()
    chunks = []
    
    for i in range(0, len(words), chunk_size - overlap):
        chunk = ' '.join(words[i:i + chunk_size])
        if chunk.strip():
            chunks.append(chunk)
        
        # Break if we're at the end
        if i + chunk_size >= len(words):
            break
    
    return chunks

def process_document(input_file, output_dir, chunk_size=200):
    """Xử lý một tài liệu thành nhiều file nhỏ"""
    
    # Đọc file input
    with open(input_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Làm sạch text
    content = re.sub(r'\n+', '\n', content)
    content = content.strip()
    
    # Chia thành chunks
    chunks = chunk_text(content, chunk_size)
    
    # Tạo thư mục output
    os.makedirs(output_dir, exist_ok=True)
    
    # Lưu từng chunk
    base_name = os.path.splitext(os.path.basename(input_file))[0]
    
    for i, chunk in enumerate(chunks):
        output_file = os.path.join(output_dir, f"{base_name}_chunk_{i+1:03d}.txt")
        
        with open(output_file, 'w', encoding='utf-8') as f:
            f.write(f"# {base_name} - Phần {i+1}/{len(chunks)}\n\n")
            f.write(chunk)
        
        print(f"Created: {output_file}")
    
    print(f"\nProcessed {input_file} into {len(chunks)} chunks")
    return len(chunks)

def main():
    if len(sys.argv) != 2:
        print("Usage: python3 chunk-document.py <input_file>")
        print("Example: python3 chunk-document.py test_documents/quy_tac_tinh_phep_nam.txt")
        sys.exit(1)
    
    input_file = sys.argv[1]
    
    if not os.path.exists(input_file):
        print(f"Error: File {input_file} not found")
        sys.exit(1)
    
    # Tạo thư mục output
    output_dir = "chunked_documents"
    
    print(f"Chunking document: {input_file}")
    print(f"Output directory: {output_dir}")
    print("-" * 50)
    
    total_chunks = process_document(input_file, output_dir, chunk_size=150)
    
    print(f"\n✅ Successfully created {total_chunks} chunks")
    print(f"📁 Files saved in: {output_dir}/")
    print("\nYou can now upload these smaller files individually!")

if __name__ == "__main__":
    main()


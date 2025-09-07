# Semantic Chunking Implementation

## Tổng quan

Semantic chunking là một kỹ thuật chia nhỏ văn bản dựa trên ý nghĩa và ngữ cảnh thay vì chỉ dựa trên kích thước cố định. Điều này giúp cải thiện chất lượng tìm kiếm và truy xuất thông tin trong hệ thống AI.

## Lợi ích của Semantic Chunking

### 1. Chia nhỏ theo ranh giới ngữ nghĩa
- Tự động nhận diện các đoạn văn có ý nghĩa hoàn chỉnh
- Giữ nguyên ngữ cảnh và mối liên hệ giữa các ý tưởng
- Tránh cắt ngang giữa các câu hoặc đoạn văn quan trọng

### 2. Cải thiện chất lượng embedding
- Mỗi chunk chứa thông tin có ý nghĩa hoàn chỉnh
- Vector embedding phản ánh chính xác hơn nội dung
- Tăng độ chính xác của semantic search

### 3. Tối ưu hóa hiệu suất
- Giảm số lượng chunks không cần thiết
- Tăng chất lượng kết quả tìm kiếm
- Giảm noise trong kết quả

## Cấu trúc Implementation

### 1. SemanticChunkingService
```go
type SemanticChunkingService struct {
    db           *gorm.DB
    geminiClient *GeminiClientV2
}
```

### 2. ChunkConfig
```go
type ChunkConfig struct {
    MinChunkSize          int     // Kích thước tối thiểu (200 chars)
    MaxChunkSize          int     // Kích thước tối đa (1000 chars)
    SimilarityThreshold   float64 // Ngưỡng tương đồng (0.7)
    OverlapSize          int     // Độ chồng lấp (100 chars)
    UseSemanticBoundaries bool    // Sử dụng ranh giới ngữ nghĩa
}
```

## Quy trình Semantic Chunking

### 1. Preprocessing
- Chuẩn hóa khoảng trắng
- Loại bỏ dấu câu thừa
- Đảm bảo kết thúc câu đúng

### 2. Identify Semantic Boundaries
- **Paragraph boundaries**: Phát hiện ranh giới đoạn văn
- **Sentence boundaries**: Tìm kết thúc câu trong đoạn dài
- **Topic boundaries**: Nhận diện tiêu đề, danh sách, headers

### 3. Create Semantic Chunks
- Chia nhỏ dựa trên ranh giới ngữ nghĩa
- Đảm bảo kích thước chunk trong khoảng cho phép
- Áp dụng overlap giữa các chunks

### 4. Generate Embeddings
- Tạo embedding cho mỗi semantic chunk
- Lưu trữ vào database với metadata

## Cách sử dụng

### 1. Upload Document (Tự động sử dụng Semantic Chunking)
```bash
curl -X POST http://localhost:8080/api/v1/documents/upload \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Document",
    "content": "Nội dung tài liệu..."
  }'
```

### 2. Re-embed với Semantic Chunking
```bash
curl -X POST http://localhost:8080/api/v1/documents/semantic-reembed \
  -H "Content-Type: application/json" \
  -d '{
    "document_id": "uuid-here",
    "config": {
      "min_chunk_size": 200,
      "max_chunk_size": 800,
      "similarity_threshold": 0.7,
      "overlap_size": 100,
      "use_semantic_boundaries": true
    }
  }'
```

### 3. Test Scripts
```bash
# Test semantic chunking
./scripts/test-semantic-chunking.sh

# So sánh kết quả chunking
python3 scripts/compare-chunking.py
```

## Cấu hình tối ưu

### Cho tài liệu ngắn (< 2000 từ)
```go
config := &ChunkConfig{
    MinChunkSize:          150,
    MaxChunkSize:          600,
    SimilarityThreshold:   0.8,
    OverlapSize:          50,
    UseSemanticBoundaries: true,
}
```

### Cho tài liệu dài (> 5000 từ)
```go
config := &ChunkConfig{
    MinChunkSize:          300,
    MaxChunkSize:          1000,
    SimilarityThreshold:   0.6,
    OverlapSize:          150,
    UseSemanticBoundaries: true,
}
```

### Cho tài liệu kỹ thuật
```go
config := &ChunkConfig{
    MinChunkSize:          400,
    MaxChunkSize:          1200,
    SimilarityThreshold:   0.7,
    OverlapSize:          200,
    UseSemanticBoundaries: true,
}
```

## So sánh với Legacy Chunking

| Aspect | Legacy Chunking | Semantic Chunking |
|--------|----------------|-------------------|
| **Phương pháp** | Chia theo kích thước cố định | Chia theo ý nghĩa |
| **Ranh giới** | Có thể cắt ngang câu | Tôn trọng ranh giới ngữ nghĩa |
| **Chất lượng** | Trung bình | Cao |
| **Hiệu suất** | Nhanh | Chậm hơn (do xử lý phức tạp) |
| **Độ chính xác** | Thấp | Cao |
| **Ngữ cảnh** | Có thể mất ngữ cảnh | Giữ nguyên ngữ cảnh |

## Monitoring và Debug

### 1. Log Messages
```
Starting semantic chunking for document: example.pdf
Split into 15 chunks
Processing semantic chunk 1/15
Generated embedding with 768 dimensions
Saved chunk 1 to database
Completed semantic chunking for document: example.pdf (15 chunks)
```

### 2. Error Handling
- Fallback tự động về legacy chunking nếu semantic chunking thất bại
- Log chi tiết các lỗi xảy ra
- Retry mechanism cho API calls

### 3. Performance Metrics
- Thời gian xử lý per document
- Số lượng chunks được tạo
- Chất lượng search results

## Tương lai và Cải tiến

### 1. Advanced Semantic Analysis
- Sử dụng transformer models để phân tích ngữ nghĩa
- Topic modeling để nhóm các chunks liên quan
- Sentiment analysis để ưu tiên chunks quan trọng

### 2. Dynamic Configuration
- Tự động điều chỉnh config dựa trên loại tài liệu
- Machine learning để tối ưu hóa parameters
- A/B testing để so sánh hiệu quả

### 3. Multi-language Support
- Hỗ trợ phân tích ngữ nghĩa cho nhiều ngôn ngữ
- Language-specific boundary detection
- Cross-language semantic similarity

## Troubleshooting

### 1. Chunks quá nhỏ
- Tăng `MinChunkSize`
- Giảm `SimilarityThreshold`
- Kiểm tra preprocessing logic

### 2. Chunks quá lớn
- Giảm `MaxChunkSize`
- Tăng `SimilarityThreshold`
- Cải thiện boundary detection

### 3. Performance chậm
- Giảm `OverlapSize`
- Tối ưu hóa API calls
- Sử dụng batch processing

### 4. Chất lượng search kém
- Điều chỉnh `SimilarityThreshold`
- Cải thiện preprocessing
- Kiểm tra embedding quality


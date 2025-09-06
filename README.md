# Company AI Training System

Hệ thống AI để training và truy vấn tài liệu công ty sử dụng Golang + Google Gemini API + PostgreSQL với pgvector.

## Tính năng

- 📄 **Upload và xử lý tài liệu**: Hỗ trợ PDF, DOCX, TXT
- 🔍 **Vector Search**: Tìm kiếm semantic dựa trên embeddings
- 💬 **Chat AI**: Trò chuyện với AI dựa trên tài liệu công ty
- 🎯 **RAG System**: Retrieval-Augmented Generation cho câu trả lời chính xác
- 🐳 **Docker Support**: Dễ dàng triển khai với Docker Compose

## Kiến trúc hệ thống

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend API   │    │   PostgreSQL    │
│   (React)       │───▶│   (Golang)      │───▶│   + pgvector    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │  Google Gemini  │
                       │     API         │
                       └─────────────────┘
```

## Cài đặt và chạy

### Quick Start (Khuyến nghị)

**Chạy cả backend và frontend với một lệnh:**

```bash
./scripts/start-dev.sh
```

Script này sẽ:
- Khởi động PostgreSQL và Go backend qua Docker
- Cài đặt frontend dependencies  
- Khởi động React development server
- Mở trình duyệt tại http://localhost:3000

**Lưu ý**: Bạn cần có Google Gemini API key để sử dụng hệ thống. Xem phần cấu hình bên dưới.

**Hoặc test riêng API:**

```bash
./scripts/test-api.sh
```

### Manual Setup

### 1. Cài đặt Dependencies

```bash
# Cài đặt Go dependencies
go mod tidy

# Hoặc sử dụng Docker (khuyến nghị)
docker-compose up -d
```

### 2. Cấu hình môi trường

```bash
# Copy file cấu hình mẫu
cp env.example .env

# Chỉnh sửa cấu hình theo nhu cầu
nano .env
```

**Quan trọng**: Bạn cần lấy Google Gemini API key:

1. Truy cập https://makersuite.google.com/app/apikey
2. Đăng nhập với tài khoản Google
3. Tạo API key mới
4. Copy API key vào file `.env`:

```bash
GEMINI_API_KEY=your_actual_api_key_here
```

### 3. Khởi động hệ thống

#### Với Docker (Khuyến nghị)

```bash
# Khởi động tất cả services
docker-compose up -d
```

#### Chạy local

```bash
# Khởi động PostgreSQL trước
docker-compose up -d postgres

# Chạy ứng dụng
go run main.go
```

### 4. Kiểm tra hệ thống

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Kết quả mong đợi:
# {"service":"company-ai-training","status":"ok"}
```

## API Endpoints

### Documents

- `POST /api/v1/documents/upload` - Upload tài liệu
- `GET /api/v1/documents` - Lấy danh sách tài liệu
- `GET /api/v1/documents/:id` - Lấy chi tiết tài liệu
- `DELETE /api/v1/documents/:id` - Xóa tài liệu

### Search

- `GET /api/v1/search?q=query&limit=10` - Tìm kiếm tài liệu

### Chat

- `POST /api/v1/chat/sessions` - Tạo phiên chat mới
- `GET /api/v1/chat/sessions` - Lấy danh sách phiên chat
- `GET /api/v1/chat/sessions/:id` - Lấy chi tiết phiên chat
- `POST /api/v1/chat/sessions/:id/messages` - Gửi tin nhắn
- `DELETE /api/v1/chat/sessions/:id` - Xóa phiên chat

## Ví dụ sử dụng

### 1. Upload tài liệu

```bash
curl -X POST \
  http://localhost:8080/api/v1/documents/upload \
  -H 'Content-Type: multipart/form-data' \
  -F 'file=@/path/to/document.pdf'
```

### 2. Tạo phiên chat

```bash
curl -X POST \
  http://localhost:8080/api/v1/chat/sessions \
  -H 'Content-Type: application/json' \
  -d '{"name": "Chat về chính sách công ty"}'
```

### 3. Gửi tin nhắn

```bash
curl -X POST \
  http://localhost:8080/api/v1/chat/sessions/{session_id}/messages \
  -H 'Content-Type: application/json' \
  -d '{"message": "Chính sách nghỉ phép của công ty như thế nào?"}'
```

### 4. Tìm kiếm tài liệu

```bash
curl "http://localhost:8080/api/v1/search?q=chính%20sách%20nhân%20sự&limit=5"
```

## Cấu trúc project

```
company-ai-training/
├── main.go                    # Entry point
├── go.mod                     # Go modules
├── go.sum                     # Go dependencies
├── Dockerfile                 # Docker build file
├── docker-compose.yml         # Docker compose config
├── env.example               # Environment config example
├── README.md                 # Documentation
└── internal/
    ├── api/                  # REST API handlers
    │   ├── handlers.go
    │   └── server.go
    ├── config/               # Configuration
    │   └── config.go
    ├── database/             # Database setup
    │   └── database.go
    ├── models/               # Data models
    │   └── document.go
    └── services/             # Business logic
        ├── document_service.go
        ├── vector_service.go
        ├── chat_service.go
        ├── user_service.go
        └── gemini_client.go
```

## Models được sử dụng

### Chat Model
- **Gemini 1.5 Flash**: Model chat chính, nhanh và hiệu quả

### Embedding Model
- **Embedding-001**: Model embedding của Google cho việc tìm kiếm semantic

## Lưu ý quan trọng

1. **API Key**: Cần Google Gemini API key để sử dụng (miễn phí với giới hạn)
2. **Internet Connection**: Cần kết nối internet để gọi Gemini API
3. **Vector Extension**: Cần PostgreSQL với pgvector extension
4. **File Types**: Chỉ hỗ trợ PDF, DOCX, TXT
5. **Rate Limits**: Gemini API có giới hạn request/phút (60 requests/minute cho free tier)

## Monitoring và Logs

```bash
# Xem logs containers
docker-compose logs -f app
docker-compose logs -f postgres

# Kiểm tra trạng thái
docker-compose ps
```

## Troubleshooting

### Lỗi thường gặp

1. **Gemini API error**: Kiểm tra API key và kết nối internet
2. **Database connection failed**: Kiểm tra PostgreSQL đã khởi động
3. **Vector extension not found**: Sử dụng image `pgvector/pgvector:pg15`
4. **Rate limit exceeded**: Đợi một chút rồi thử lại (60 requests/minute limit)

### Debug

```bash
# Test Gemini API connection
curl -H "Content-Type: application/json" \
     -d '{"contents":[{"parts":[{"text":"Hello"}]}]}' \
     "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=YOUR_API_KEY"

# Kiểm tra database
docker exec -it company_ai_postgres psql -U postgres -d company_ai -c "SELECT * FROM documents LIMIT 5;"

# Kiểm tra environment variables
echo $GEMINI_API_KEY
```

## Phát triển thêm

- [ ] Thêm authentication/authorization
- [ ] Web UI cho admin
- [ ] Hỗ trợ thêm file formats (Excel, PowerPoint)
- [ ] Real-time chat với WebSocket
- [ ] Caching layer (Redis)
- [ ] Metrics và monitoring
- [ ] Multi-language support

## License

MIT License

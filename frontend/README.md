# Company AI Chat Frontend

Ứng dụng chat frontend được xây dựng với ReactJS để tương tác với Company AI Assistant.

## Tính năng

- 💬 **Chat thông minh**: Trò chuyện với AI assistant về chính sách công ty
- 📚 **Tra cứu tài liệu**: Tìm kiếm thông tin từ các tài liệu nội bộ
- 🔄 **Quản lý phiên chat**: Tạo, xóa và chuyển đổi giữa các cuộc trò chuyện
- 📱 **Responsive**: Tối ưu cho cả desktop và mobile
- 🎨 **UI hiện đại**: Giao diện đẹp mắt và dễ sử dụng

## Cài đặt

### Yêu cầu hệ thống
- Node.js >= 14.0.0
- npm hoặc yarn

### Cài đặt dependencies

```bash
cd frontend
npm install
```

### Cấu hình môi trường

Sao chép file `env.example` thành `.env`:

```bash
cp env.example .env
```

Chỉnh sửa file `.env` với URL backend của bạn:

```env
REACT_APP_API_URL=http://localhost:8082/api/v1
```

## Chạy ứng dụng

### Development mode

```bash
npm start
```

Ứng dụng sẽ chạy tại `http://localhost:3000`

### Production build

```bash
npm run build
```

## Cấu trúc thư mục

```
src/
├── components/           # React components
│   ├── ChatMessage.js   # Component hiển thị tin nhắn
│   ├── ChatInput.js     # Component nhập tin nhắn
│   ├── ChatWindow.js    # Component cửa sổ chat chính
│   └── SessionList.js   # Component danh sách phiên chat
├── services/
│   └── api.js          # API client
├── App.js              # Component chính
├── App.css             # Styles chính
└── index.js            # Entry point
```

## API Integration

Ứng dụng kết nối với backend API qua các endpoints:

- `POST /api/v1/chat/sessions` - Tạo phiên chat mới
- `GET /api/v1/chat/sessions` - Lấy danh sách phiên chat
- `GET /api/v1/chat/sessions/:id` - Lấy chi tiết phiên chat
- `POST /api/v1/chat/sessions/:id/messages` - Gửi tin nhắn
- `DELETE /api/v1/chat/sessions/:id` - Xóa phiên chat

## Tính năng chính

### 1. Quản lý phiên chat
- Tạo phiên chat mới
- Xem danh sách tất cả phiên chat
- Chuyển đổi giữa các phiên chat
- Xóa phiên chat không cần thiết

### 2. Giao diện chat
- Hiển thị tin nhắn theo thời gian thực
- Phân biệt tin nhắn người dùng và AI
- Indicator khi AI đang trả lời
- Cuộn tự động đến tin nhắn mới nhất

### 3. Responsive design
- Tối ưu cho desktop (>768px)
- Sidebar ẩn/hiện trên mobile
- Touch-friendly interface
- Adaptive layouts

## Customization

### Thay đổi theme colors

Chỉnh sửa CSS variables trong `src/index.css`:

```css
:root {
  --primary-color: #667eea;
  --secondary-color: #764ba2;
  --background-color: #f8f9fa;
  --text-color: #333;
}
```

### Thêm tính năng mới

1. Tạo component mới trong `src/components/`
2. Thêm API calls trong `src/services/api.js`
3. Import và sử dụng trong `App.js`

## Troubleshooting

### Lỗi kết nối API
- Kiểm tra backend server đã chạy
- Xác minh `REACT_APP_API_URL` trong `.env`
- Kiểm tra CORS settings trên backend

### Lỗi build
- Xóa `node_modules` và chạy lại `npm install`
- Kiểm tra Node.js version
- Xóa cache: `npm start -- --reset-cache`

## Browser Support

- Chrome >= 60
- Firefox >= 60
- Safari >= 12
- Edge >= 79

## Performance

- Code splitting với React.lazy (có thể thêm)
- Optimized images và assets
- Efficient re-rendering với React.memo
- Lazy loading cho large lists

## Security

- Input sanitization
- XSS protection
- HTTPS trong production
- Environment variables cho sensitive data

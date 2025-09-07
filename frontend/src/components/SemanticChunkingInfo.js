import React from 'react';
import { Brain, Info, CheckCircle, AlertCircle } from 'lucide-react';
import './SemanticChunkingInfo.css';

const SemanticChunkingInfo = ({ isOpen, onClose }) => {
  if (!isOpen) return null;

  return (
    <div className="semantic-info-overlay">
      <div className="semantic-info-modal">
        <div className="info-header">
          <div className="info-title">
            <Brain size={20} />
            <h3>Semantic Chunking là gì?</h3>
          </div>
          <button onClick={onClose} className="close-btn">×</button>
        </div>

        <div className="info-content">
          <div className="info-section">
            <h4>🎯 Khái niệm</h4>
            <p>
              Semantic chunking là kỹ thuật chia nhỏ văn bản dựa trên ý nghĩa và ngữ cảnh 
              thay vì chỉ dựa trên kích thước cố định. Điều này giúp cải thiện chất lượng 
              tìm kiếm và truy xuất thông tin.
            </p>
          </div>

          <div className="info-section">
            <h4>✅ Lợi ích</h4>
            <ul>
              <li>
                <CheckCircle size={16} />
                <span>Chia nhỏ theo ranh giới ngữ nghĩa tự nhiên</span>
              </li>
              <li>
                <CheckCircle size={16} />
                <span>Giữ nguyên ngữ cảnh và ý nghĩa của từng đoạn</span>
              </li>
              <li>
                <CheckCircle size={16} />
                <span>Cải thiện độ chính xác của vector search</span>
              </li>
              <li>
                <CheckCircle size={16} />
                <span>Tối ưu hóa hiệu suất embedding</span>
              </li>
            </ul>
          </div>

          <div className="info-section">
            <h4>⚙️ Cách hoạt động</h4>
            <ol>
              <li><strong>Preprocessing:</strong> Chuẩn hóa và làm sạch văn bản</li>
              <li><strong>Boundary Detection:</strong> Tìm ranh giới đoạn văn, câu, và chủ đề</li>
              <li><strong>Semantic Analysis:</strong> Phân tích ý nghĩa để tạo chunks có nghĩa</li>
              <li><strong>Size Optimization:</strong> Đảm bảo kích thước chunk phù hợp</li>
              <li><strong>Overlap Application:</strong> Thêm độ chồng lấp giữa các chunks</li>
            </ol>
          </div>

          <div className="info-section">
            <h4>📊 So sánh với Legacy Chunking</h4>
            <div className="comparison-table">
              <div className="comparison-row header">
                <div>Tiêu chí</div>
                <div>Legacy Chunking</div>
                <div>Semantic Chunking</div>
              </div>
              <div className="comparison-row">
                <div>Phương pháp</div>
                <div>Chia theo kích thước cố định</div>
                <div>Chia theo ý nghĩa</div>
              </div>
              <div className="comparison-row">
                <div>Ranh giới</div>
                <div>Có thể cắt ngang câu</div>
                <div>Tôn trọng ranh giới ngữ nghĩa</div>
              </div>
              <div className="comparison-row">
                <div>Chất lượng</div>
                <div>Trung bình</div>
                <div>Cao</div>
              </div>
              <div className="comparison-row">
                <div>Độ chính xác</div>
                <div>Thấp</div>
                <div>Cao</div>
              </div>
            </div>
          </div>

          <div className="info-section">
            <h4>💡 Khuyến nghị sử dụng</h4>
            <div className="recommendations">
              <div className="recommendation-item">
                <AlertCircle size={16} className="warning-icon" />
                <div>
                  <strong>Tài liệu ngắn:</strong> Sử dụng preset "Tài liệu ngắn" với 
                  kích thước chunk nhỏ hơn
                </div>
              </div>
              <div className="recommendation-item">
                <AlertCircle size={16} className="warning-icon" />
                <div>
                  <strong>Tài liệu dài:</strong> Sử dụng preset "Tài liệu dài" với 
                  kích thước chunk lớn hơn
                </div>
              </div>
              <div className="recommendation-item">
                <AlertCircle size={16} className="warning-icon" />
                <div>
                  <strong>Tài liệu kỹ thuật:</strong> Sử dụng preset "Tài liệu kỹ thuật" 
                  với độ chồng lấp cao hơn
                </div>
              </div>
            </div>
          </div>
        </div>

        <div className="info-actions">
          <button onClick={onClose} className="close-info-btn">
            Đóng
          </button>
        </div>
      </div>
    </div>
  );
};

export default SemanticChunkingInfo;


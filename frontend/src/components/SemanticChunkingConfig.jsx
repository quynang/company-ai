import React, { useState } from 'react';
import { Settings, Brain, Save, X, Info } from 'lucide-react';
import SemanticChunkingInfo from './SemanticChunkingInfo';
import './SemanticChunkingConfig.css';

const SemanticChunkingConfig = ({ isOpen, onClose, onSave, initialConfig = null }) => {
  const [config, setConfig] = useState({
    minChunkSize: initialConfig?.minChunkSize || 200,
    maxChunkSize: initialConfig?.maxChunkSize || 1000,
    similarityThreshold: initialConfig?.similarityThreshold || 0.7,
    overlapSize: initialConfig?.overlapSize || 100,
    useSemanticBoundaries: initialConfig?.useSemanticBoundaries !== false
  });

  const [preset, setPreset] = useState('default');
  const [showInfo, setShowInfo] = useState(false);

  const presets = {
    default: {
      name: 'Mặc định',
      config: { minChunkSize: 200, maxChunkSize: 1000, similarityThreshold: 0.7, overlapSize: 100, useSemanticBoundaries: true }
    },
    short: {
      name: 'Tài liệu ngắn',
      config: { minChunkSize: 150, maxChunkSize: 600, similarityThreshold: 0.8, overlapSize: 50, useSemanticBoundaries: true }
    },
    long: {
      name: 'Tài liệu dài',
      config: { minChunkSize: 300, maxChunkSize: 1000, similarityThreshold: 0.6, overlapSize: 150, useSemanticBoundaries: true }
    },
    technical: {
      name: 'Tài liệu kỹ thuật',
      config: { minChunkSize: 400, maxChunkSize: 1200, similarityThreshold: 0.7, overlapSize: 200, useSemanticBoundaries: true }
    }
  };

  const handlePresetChange = (presetKey) => {
    setPreset(presetKey);
    setConfig(presets[presetKey].config);
  };

  const handleInputChange = (field, value) => {
    setConfig(prev => ({
      ...prev,
      [field]: field === 'useSemanticBoundaries' ? value : 
               field === 'similarityThreshold' ? parseFloat(value) : parseInt(value)
    }));
  };

  const handleSave = () => {
    onSave(config);
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="semantic-config-overlay">
      <div className="semantic-config-modal">
        <div className="config-header">
          <div className="config-title">
            <Brain size={20} />
            <h3>Cấu hình Semantic Chunking</h3>
          </div>
          <div className="header-actions">
            <button 
              onClick={() => setShowInfo(true)} 
              className="info-btn"
              title="Tìm hiểu về Semantic Chunking"
            >
              <Info size={16} />
            </button>
            <button onClick={onClose} className="close-btn">
              <X size={20} />
            </button>
          </div>
        </div>

        <div className="config-content">
          <div className="preset-section">
            <label>Chọn preset:</label>
            <div className="preset-buttons">
              {Object.entries(presets).map(([key, preset]) => (
                <button
                  key={key}
                  className={`preset-btn ${preset === key ? 'active' : ''}`}
                  onClick={() => handlePresetChange(key)}
                >
                  {preset.name}
                </button>
              ))}
            </div>
          </div>

          <div className="config-section">
            <h4>Thông số cấu hình</h4>
            
            <div className="config-grid">
              <div className="config-item">
                <label htmlFor="minChunkSize">Kích thước tối thiểu (ký tự):</label>
                <input
                  id="minChunkSize"
                  type="number"
                  min="50"
                  max="1000"
                  value={config.minChunkSize}
                  onChange={(e) => handleInputChange('minChunkSize', e.target.value)}
                />
                <small>Kích thước nhỏ nhất cho mỗi chunk</small>
              </div>

              <div className="config-item">
                <label htmlFor="maxChunkSize">Kích thước tối đa (ký tự):</label>
                <input
                  id="maxChunkSize"
                  type="number"
                  min="200"
                  max="2000"
                  value={config.maxChunkSize}
                  onChange={(e) => handleInputChange('maxChunkSize', e.target.value)}
                />
                <small>Kích thước lớn nhất cho mỗi chunk</small>
              </div>

              <div className="config-item">
                <label htmlFor="similarityThreshold">Ngưỡng tương đồng (0-1):</label>
                <input
                  id="similarityThreshold"
                  type="number"
                  min="0"
                  max="1"
                  step="0.1"
                  value={config.similarityThreshold}
                  onChange={(e) => handleInputChange('similarityThreshold', e.target.value)}
                />
                <small>Ngưỡng để xác định ranh giới ngữ nghĩa</small>
              </div>

              <div className="config-item">
                <label htmlFor="overlapSize">Độ chồng lấp (ký tự):</label>
                <input
                  id="overlapSize"
                  type="number"
                  min="0"
                  max="500"
                  value={config.overlapSize}
                  onChange={(e) => handleInputChange('overlapSize', e.target.value)}
                />
                <small>Số ký tự chồng lấp giữa các chunks</small>
              </div>
            </div>

            <div className="config-item checkbox-item">
              <label className="checkbox-label">
                <input
                  type="checkbox"
                  checked={config.useSemanticBoundaries}
                  onChange={(e) => handleInputChange('useSemanticBoundaries', e.target.checked)}
                />
                <span className="checkmark"></span>
                Sử dụng ranh giới ngữ nghĩa
              </label>
              <small>Chia nhỏ dựa trên ý nghĩa thay vì kích thước cố định</small>
            </div>
          </div>

          <div className="config-info">
            <h4>Thông tin về Semantic Chunking</h4>
            <ul>
              <li><strong>Kích thước tối thiểu:</strong> Đảm bảo chunks không quá nhỏ</li>
              <li><strong>Kích thước tối đa:</strong> Tránh chunks quá lớn</li>
              <li><strong>Ngưỡng tương đồng:</strong> Cao hơn = chunks nhỏ hơn, thấp hơn = chunks lớn hơn</li>
              <li><strong>Độ chồng lấp:</strong> Giúp giữ ngữ cảnh giữa các chunks</li>
              <li><strong>Ranh giới ngữ nghĩa:</strong> Chia nhỏ theo ý nghĩa thay vì kích thước</li>
            </ul>
          </div>
        </div>

        <div className="config-actions">
          <button onClick={onClose} className="cancel-btn">
            Hủy
          </button>
          <button onClick={handleSave} className="save-btn">
            <Save size={16} />
            Lưu cấu hình
          </button>
        </div>
      </div>

      {/* Info Modal */}
      <SemanticChunkingInfo
        isOpen={showInfo}
        onClose={() => setShowInfo(false)}
      />
    </div>
  );
};

export default SemanticChunkingConfig;

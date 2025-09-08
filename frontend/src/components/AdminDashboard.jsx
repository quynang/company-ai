import React, { useState, useEffect } from 'react';
import { Edit, RefreshCw, Trash2, FileText, Plus, Search, Save, X, Brain, Settings, Tag } from 'lucide-react';
import ConfirmDialog from './ConfirmDialog';
import SemanticChunkingConfig from './SemanticChunkingConfig';
import CategoryManager from './CategoryManager';
import CategoryMultiSelect from './CategoryMultiSelect';
import { documentAPI } from '../services/api';
import './AdminDashboard.css';

const AdminDashboard = () => {
  const [documents, setDocuments] = useState([]);
  const [loading, setLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedDoc, setSelectedDoc] = useState(null);
  const [editingDoc, setEditingDoc] = useState(null);
  const [editContent, setEditContent] = useState('');
  const [newDocName, setNewDocName] = useState('');
  const [newDocContent, setNewDocContent] = useState('');
  const [newDocCategories, setNewDocCategories] = useState([]);
  const [isCreating, setIsCreating] = useState(false);
  const [selectedDocCategories, setSelectedDocCategories] = useState([]);
  const [isUpdatingCategories, setIsUpdatingCategories] = useState(false);
  const [confirmDialog, setConfirmDialog] = useState({
    isOpen: false,
    title: '',
    message: '',
    onConfirm: null,
    type: 'danger'
  });
  const [semanticConfigOpen, setSemanticConfigOpen] = useState(false);
  const [semanticConfig, setSemanticConfig] = useState(null);
  const [activeTab, setActiveTab] = useState('documents'); // 'documents' or 'categories'

  useEffect(() => {
    fetchDocuments();
  }, []);

  const fetchDocuments = async () => {
    setLoading(true);
    try {
      const data = await documentAPI.getDocuments();
      setDocuments(data.documents || []);
    } catch (error) {
      console.error('Error fetching documents:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadDocumentCategories = async (documentId) => {
    try {
      const response = await documentAPI.getDocumentCategories(documentId);
      setSelectedDocCategories(response.categories?.map(cat => cat.id) || []);
    } catch (error) {
      console.error('Error loading document categories:', error);
      setSelectedDocCategories([]);
    }
  };

  const handleUpdateDocumentCategories = async () => {
    if (!selectedDoc) return;
    
    setIsUpdatingCategories(true);
    try {
      await documentAPI.updateDocumentCategories(selectedDoc.id, selectedDocCategories);
      alert('Document categories updated successfully!');
      // Reload documents to refresh the data
      fetchDocuments();
    } catch (error) {
      console.error('Error updating document categories:', error);
      alert('Failed to update document categories');
    } finally {
      setIsUpdatingCategories(false);
    }
  };

  const handleCreateNew = async () => {
    if (!newDocName.trim() || !newDocContent.trim()) {
      alert('Vui lòng nhập tên và nội dung document');
      return;
    }

    try {
      // Create a new document with content and categories
      const result = await documentAPI.uploadDocumentFromText(newDocName, newDocContent, newDocCategories);
      alert('Tạo document thành công! Document đang được xử lý embedding...');
      setNewDocName('');
      setNewDocContent('');
      setNewDocCategories([]);
      setIsCreating(false);
      fetchDocuments();
    } catch (error) {
      console.error('Create error:', error);
      alert('Có lỗi xảy ra khi tạo document');
    }
  };

  const handleUpdate = async (docId) => {
    if (!editContent.trim()) {
      alert('Vui lòng nhập nội dung');
      return;
    }

    try {
      await documentAPI.updateDocument(docId, editContent);
      alert('Cập nhật thành công! Document đang được xử lý embedding lại...');
      setEditingDoc(null);
      setEditContent('');
      fetchDocuments();
    } catch (error) {
      console.error('Update error:', error);
      alert('Có lỗi xảy ra khi cập nhật');
    }
  };

  const handleDelete = async (docId) => {
    setConfirmDialog({
      isOpen: true,
      title: 'Xác nhận xóa',
      message: 'Bạn có chắc muốn xóa document này?',
      onConfirm: () => performDelete(docId),
      type: 'danger'
    });
  };

  const performDelete = async (docId) => {
    setConfirmDialog({ isOpen: false, title: '', message: '', onConfirm: null, type: 'danger' });

    try {
      await documentAPI.deleteDocument(docId);
      alert('Xóa thành công!');
      fetchDocuments();
    } catch (error) {
      console.error('Delete error:', error);
      alert('Có lỗi xảy ra khi xóa');
    }
  };

  const handleReembed = async (docId) => {
    setConfirmDialog({
      isOpen: true,
      title: 'Xác nhận Re-embedding',
      message: 'Bạn có chắc muốn chạy lại embedding cho document này?',
      onConfirm: () => performReembed(docId),
      type: 'warning'
    });
  };

  const performReembed = async (docId) => {
    setConfirmDialog({ isOpen: false, title: '', message: '', onConfirm: null, type: 'danger' });

    try {
      await documentAPI.reembedDocument(docId);
      alert('Đang chạy lại embedding...');
      fetchDocuments();
    } catch (error) {
      console.error('Re-embedding error:', error);
      alert('Có lỗi xảy ra khi chạy lại embedding');
    }
  };

  const handleSemanticReembed = async (docId) => {
    setConfirmDialog({
      isOpen: true,
      title: 'Xác nhận Semantic Re-embedding',
      message: 'Bạn có chắc muốn chạy lại embedding với semantic chunking cho document này?',
      onConfirm: () => performSemanticReembed(docId),
      type: 'warning'
    });
  };

  const performSemanticReembed = async (docId) => {
    setConfirmDialog({ isOpen: false, title: '', message: '', onConfirm: null, type: 'danger' });

    try {
      await documentAPI.reembedWithSemanticChunking(docId, semanticConfig);
      alert('Đang chạy lại embedding với semantic chunking...');
      fetchDocuments();
    } catch (error) {
      console.error('Semantic re-embedding error:', error);
      alert('Có lỗi xảy ra khi chạy lại semantic embedding');
    }
  };

  const handleSemanticConfigSave = (config) => {
    setSemanticConfig(config);
    alert('Cấu hình semantic chunking đã được lưu!');
  };

  const filteredDocuments = documents.filter(doc =>
    doc.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    doc.content.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const closeConfirmDialog = () => {
    setConfirmDialog({ isOpen: false, title: '', message: '', onConfirm: null, type: 'danger' });
  };

  return (
    <div className="admin-dashboard">
      <div className="dashboard-container">
        <div className="dashboard-header">
          <div>
            <h1>Admin Dashboard</h1>
          </div>
          <div className="header-actions">
            <button 
              onClick={() => setSemanticConfigOpen(true)} 
              className="config-btn"
              title="Cấu hình Semantic Chunking"
            >
              <Settings size={16} />
              Cấu hình Semantic
            </button>
            <button 
              onClick={() => setIsCreating(true)} 
              className="create-new-btn"
            >
              <Plus size={16} />
              Tạo Document Mới
            </button>
          </div>
        </div>

      {/* Tab Navigation */}
      <div className="tab-navigation">
        <button
          className={`tab-button ${activeTab === 'documents' ? 'active' : ''}`}
          onClick={() => setActiveTab('documents')}
        >
          <FileText size={14} />
          Documents
        </button>
        <button
          className={`tab-button ${activeTab === 'categories' ? 'active' : ''}`}
          onClick={() => setActiveTab('categories')}
        >
          <Tag size={14} />
          Categories
        </button>
      </div>

      {activeTab === 'documents' ? (
        <div className="dashboard-layout">
          {/* Sidebar - Documents List */}
          <div className="sidebar">
          <div className="search-section">
            <div className="search-box">
              <Search size={20} />
              <input
                type="text"
                placeholder="Tìm kiếm documents..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
              />
            </div>
          </div>

          <div className="documents-list">
            {loading ? (
              <div className="loading">Đang tải...</div>
            ) : (
              filteredDocuments.map((doc) => (
                <div 
                  key={doc.id} 
                  className={`document-item ${selectedDoc?.id === doc.id ? 'active' : ''}`}
                  onClick={() => {
                    setSelectedDoc(doc);
                    loadDocumentCategories(doc.id);
                  }}
                >
                  <div className="document-icon">
                    <FileText size={16} />
                  </div>
                  <div className="document-info">
                    <h4 className="document-name">{doc.name}</h4>
                    <p className="document-description">
                      {doc.content.substring(0, 100)}...
                    </p>
                    {doc.categories && doc.categories.length > 0 && (
                      <div className="document-categories">
                        {doc.categories.map((category) => (
                          <span key={category.id}>
                            {category.name}
                          </span>
                        ))}
                      </div>
                    )}
                    <span className="document-date">
                      {new Date(doc.created_at).toLocaleDateString('vi-VN')}
                    </span>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>

        {/* Main Content - Document Editor */}
        <div className="main-content">
          {isCreating ? (
            <div className="create-document">
              <div className="create-header">
                <h2>Tạo Document Mới</h2>
                <button 
                  onClick={() => setIsCreating(false)} 
                  className="close-btn"
                >
                  <X size={20} />
                </button>
              </div>
              <div className="create-form">
                <div className="form-group">
                  <label>Tên Document:</label>
                  <input
                    type="text"
                    value={newDocName}
                    onChange={(e) => setNewDocName(e.target.value)}
                    placeholder="Nhập tên document"
                  />
                </div>
                <div className="form-group">
                  <CategoryMultiSelect
                    selectedCategoryIds={newDocCategories}
                    onCategoryChange={setNewDocCategories}
                  />
                </div>
                <div className="form-group">
                  <label>Nội dung:</label>
                  <textarea
                    value={newDocContent}
                    onChange={(e) => setNewDocContent(e.target.value)}
                    rows={20}
                    placeholder="Nhập nội dung document..."
                  />
                </div>
                <div className="form-actions">
                  <button onClick={handleCreateNew} className="save-btn">
                    <Save size={16} />
                    Tạo Document
                  </button>
                  <button 
                    onClick={() => setIsCreating(false)} 
                    className="cancel-btn"
                  >
                    Hủy
                  </button>
                </div>
              </div>
            </div>
          ) : selectedDoc ? (
            <div className="document-editor">
              <div className="editor-header">
                <h2>{selectedDoc.name}</h2>
                <div className="editor-actions">
                  <button
                    onClick={() => {
                      setEditingDoc(selectedDoc.id);
                      setEditContent(selectedDoc.content);
                    }}
                    className="action-btn edit-btn"
                    title="Chỉnh sửa"
                  >
                    <Edit size={12} />
                    Chỉnh sửa
                  </button>
                  <button
                    onClick={() => handleReembed(selectedDoc.id)}
                    className="action-btn reembed-btn"
                    title="Chạy lại embedding (Legacy)"
                  >
                    <RefreshCw size={12} />
                    Re-embed
                  </button>
                  <button
                    onClick={() => handleSemanticReembed(selectedDoc.id)}
                    className="action-btn semantic-btn"
                    title="Chạy lại embedding với Semantic Chunking"
                  >
                    <Brain size={12} />
                    Semantic Re-embed
                  </button>
                  <button
                    onClick={() => handleDelete(selectedDoc.id)}
                    className="action-btn delete-btn"
                    title="Xóa"
                  >
                    <Trash2 size={12} />
                    Xóa
                  </button>
                </div>
              </div>

              <div className="editor-content">
                {/* Category Management Section */}
                <div className="category-management">
                  <h3>
                    <Tag size={20} />
                    Document Categories
                  </h3>
                  <div className="mb-4">
                    <CategoryMultiSelect
                      selectedCategoryIds={selectedDocCategories}
                      onCategoryChange={setSelectedDocCategories}
                    />
                  </div>
                  <div className="flex">
                    <button
                      onClick={handleUpdateDocumentCategories}
                      disabled={isUpdatingCategories}
                      className="btn-update-categories"
                    >
                      <Save size={10} />
                      {isUpdatingCategories ? 'Updating...' : 'Update Categories'}
                    </button>
                  </div>
                </div>

                {editingDoc === selectedDoc.id ? (
                  <div className="edit-mode">
                    <textarea
                      value={editContent}
                      onChange={(e) => setEditContent(e.target.value)}
                      rows={25}
                      placeholder="Nhập nội dung document..."
                    />
                    <div className="edit-actions">
                      <button
                        onClick={() => handleUpdate(selectedDoc.id)}
                        className="save-btn"
                      >
                        <Save size={16} />
                        Lưu
                      </button>
                      <button
                        onClick={() => {
                          setEditingDoc(null);
                          setEditContent('');
                        }}
                        className="cancel-btn"
                      >
                        Hủy
                      </button>
                    </div>
                  </div>
                ) : (
                  <div className="document-view">
                    <div className="document-meta">
                      <span className="chunks-count">
                        Chunks: {selectedDoc.chunks_count || 0}
                      </span>
                      <span className="document-size">
                        Kích thước: {selectedDoc.content.length} ký tự
                      </span>
                    </div>
                    <div className="document-text">
                      <pre>{selectedDoc.content}</pre>
                    </div>
                  </div>
                )}
              </div>
            </div>
          ) : (
            <div className="no-selection">
              <FileText size={64} />
              <h3>Chọn một document để xem và chỉnh sửa</h3>
              <p>Hoặc tạo document mới từ sidebar</p>
            </div>
          )}
        </div>
        </div>
      ) : (
        <CategoryManager />
      )}
      </div>

      {/* Confirm Dialog */}
      <ConfirmDialog
        isOpen={confirmDialog.isOpen}
        title={confirmDialog.title}
        message={confirmDialog.message}
        onConfirm={confirmDialog.onConfirm}
        onCancel={closeConfirmDialog}
        type={confirmDialog.type}
      />

      {/* Semantic Chunking Config Modal */}
      <SemanticChunkingConfig
        isOpen={semanticConfigOpen}
        onClose={() => setSemanticConfigOpen(false)}
        onSave={handleSemanticConfigSave}
        initialConfig={semanticConfig}
      />
    </div>
  );
};

export default AdminDashboard;

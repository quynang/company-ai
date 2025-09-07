import React from 'react';
import { format } from 'date-fns';
import { MessageCircle, Plus, Trash2, Settings } from 'lucide-react';
import './SessionList.css';

const SessionList = ({ 
  sessions, 
  currentSession, 
  onSelectSession, 
  onNewSession, 
  onDeleteSession,
  onSwitchToAdmin,
  isLoading = false
}) => {
  const formatSessionTime = (dateString) => {
    try {
      const date = new Date(dateString);
      return format(date, 'dd/MM HH:mm');
    } catch (error) {
      return '';
    }
  };

  const handleDeleteClick = (e, sessionId) => {
    e.stopPropagation();
    if (window.confirm('Bạn có chắc chắn muốn xóa cuộc trò chuyện này?')) {
      onDeleteSession(sessionId);
    }
  };

  return (
    <div className="session-list">
      <div className="session-list-header">
        <h2>Cuộc trò chuyện</h2>
        <div className="header-actions">
          <button 
            className="new-session-button"
            onClick={onNewSession}
            disabled={isLoading}
            title="Tạo cuộc trò chuyện mới"
          >
            <Plus size={20} />
          </button>
          <button 
            className="admin-button"
            onClick={onSwitchToAdmin}
            title="Admin Dashboard"
          >
            <Settings size={20} />
          </button>
        </div>
      </div>


      <div className="sessions-container">
        {isLoading ? (
          <div className="loading-sessions">
            <div className="loading-spinner"></div>
            <span>Đang tải...</span>
          </div>
        ) : sessions.length === 0 ? (
          <div className="empty-sessions">
            <MessageCircle size={48} className="empty-icon" />
            <p>Chưa có cuộc trò chuyện nào</p>
            <button className="create-first-session" onClick={onNewSession}>
              Tạo cuộc trò chuyện đầu tiên
            </button>
          </div>
        ) : (
          <div className="sessions-list">
            {sessions.map((session) => (
              <div
                key={session.id}
                className={`session-item ${currentSession?.id === session.id ? 'active' : ''}`}
                onClick={() => onSelectSession(session)}
              >
                <div className="session-content">
                  <div className="session-name">
                    {session.name || 'Cuộc trò chuyện mới'}
                  </div>
                  <div className="session-time">
                    {formatSessionTime(session.updated_at)}
                  </div>
                </div>
                <button
                  className="delete-session-button"
                  onClick={(e) => handleDeleteClick(e, session.id)}
                  title="Xóa cuộc trò chuyện"
                >
                  <Trash2 size={16} />
                </button>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default SessionList;

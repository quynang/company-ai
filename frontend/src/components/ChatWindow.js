import React, { useState, useEffect, useRef } from 'react';
import { Menu, Bot, Settings } from 'lucide-react';
import ChatMessage from './ChatMessage';
import ChatInput from './ChatInput';
import './ChatWindow.css';

const ChatWindow = ({ 
  session, 
  messages, 
  onSendMessage, 
  isLoading = false,
  onToggleSidebar,
  onSwitchToAdmin,
  setMessages
}) => {
  const [isTyping, setIsTyping] = useState(false);
  const messagesEndRef = useRef(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages, isTyping]);

  const handleSendMessage = async (message) => {
    setIsTyping(true);
    try {
      await onSendMessage(message);
    } finally {
      setIsTyping(false);
    }
  };

  return (
    <div className="chat-window">
      <div className="chat-header">
        <button 
          className="sidebar-toggle"
          onClick={onToggleSidebar}
          title="Hiển thị/Ẩn danh sách cuộc trò chuyện"
        >
          <Menu size={20} />
        </button>
        <div className="chat-title">
          <Bot size={24} className="chat-icon" />
          <div className="title-text">
            <h1>Company AI Assistant</h1>
            <p>{session ? session.name : 'Chọn hoặc tạo cuộc trò chuyện mới'}</p>
          </div>
          <button 
            className="admin-toggle"
            onClick={onSwitchToAdmin}
            title="Chuyển đến Admin Dashboard"
          >
            <Settings size={20} />
          </button>
        </div>
      </div>

      <div className="chat-messages">
        {!session ? (
          <div className="welcome-message">
            <Bot size={64} className="welcome-icon" />
            <h2>Chào mừng đến với Company AI Assistant!</h2>
            <p>Tôi là trợ lý AI thông minh, sẵn sàng giúp bạn giải đáp thắc mắc về các chính sách và quy định của công ty.</p>
            <div className="welcome-features">
              <div className="feature">
                <span className="feature-icon">📚</span>
                <span>Tra cứu chính sách công ty</span>
              </div>
              <div className="feature">
                <span className="feature-icon">📝</span>
                <span>Hướng dẫn quy trình làm việc</span>
              </div>
              <div className="feature">
                <span className="feature-icon">🏖️</span>
                <span>Tính toán ngày phép</span>
              </div>
            </div>
          </div>
        ) : messages.length === 0 ? (
          <div className="empty-chat">
            <Bot size={48} className="empty-icon" />
            <h3>Bắt đầu cuộc trò chuyện</h3>
            <p>Hãy đặt câu hỏi đầu tiên để bắt đầu!</p>
          </div>
        ) : (
          <div className="messages-list">
            {messages.map((message, index) => (
              <ChatMessage 
                key={message.id || index} 
                message={message} 
                onActionClick={(result) => {
                  if (result.message) {
                    // Show success message
                    const successMessage = {
                      id: Date.now(),
                      role: 'assistant',
                      content: result.message,
                      created_at: new Date().toISOString()
                    };
                    setMessages(prev => [...prev, successMessage]);
                  } else if (result.error) {
                    // Show error message
                    const errorMessage = {
                      id: Date.now(),
                      role: 'assistant',
                      content: result.error,
                      created_at: new Date().toISOString()
                    };
                    setMessages(prev => [...prev, errorMessage]);
                  }
                }}
              />
            ))}
            {isTyping && <ChatMessage isLoading={true} />}
            <div ref={messagesEndRef} />
          </div>
        )}
      </div>

      <ChatInput
        onSendMessage={handleSendMessage}
        disabled={!session || isLoading || isTyping}
        placeholder={
          !session 
            ? "Vui lòng chọn hoặc tạo cuộc trò chuyện mới..."
            : "Nhập câu hỏi của bạn..."
        }
      />
    </div>
  );
};

export default ChatWindow;

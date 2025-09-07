import React from 'react';
import { format } from 'date-fns';
import { Bot, User, HelpCircle, Send } from 'lucide-react';
import MarkdownRenderer from './MarkdownRenderer';
import './ChatMessage.css';

const ChatMessage = ({ message, isLoading = false, onActionClick }) => {
  // Debug logging
  console.log('ChatMessage received:', message);
  console.log('Message structure:', {
    hasMessage: !!message,
    role: message?.role,
    content: message?.content,
    contentLength: message?.content?.length,
    hasActionCard: !!message?.action_card
  });
  
  const isUser = message?.role === 'user';
  const isAssistant = message?.role === 'assistant';

  const handleActionClick = async (action) => {
    if (onActionClick) {
      try {
        const response = await fetch(action.endpoint, {
          method: action.method,
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(action.payload)
        });
        
        const result = await response.json();
        onActionClick(result);
      } catch (error) {
        console.error('Action failed:', error);
        onActionClick({ error: 'Có lỗi xảy ra khi tạo ticket' });
      }
    }
  };

  if (isLoading) {
    return (
      <div className="message-container assistant">
        <div className="message-avatar">
          <Bot size={20} />
        </div>
        <div className="message-bubble">
          <div className="typing-indicator">
            <span></span>
            <span></span>
            <span></span>
          </div>
        </div>
      </div>
    );
  }

  if (!message) return null;

  return (
    <div className={`message-container ${isUser ? 'user' : 'assistant'}`}>
      <div className="message-avatar">
        {isUser ? <User size={20} /> : <Bot size={20} />}
      </div>
      <div className="message-content">
        <div className="message-bubble">
          <div className="message-text">
            {isAssistant ? (
              <MarkdownRenderer content={message.content} />
            ) : (
              message.content
            )}
          </div>
        </div>
        {message.action_card && (
          <div className="action-card">
            <div className="action-card-header">
              <HelpCircle size={16} />
              <span className="action-card-title">{message.action_card.title}</span>
            </div>
            <div className="action-card-description">
              {message.action_card.description}
            </div>
            <button 
              className="action-button"
              onClick={() => handleActionClick(message.action_card.action)}
            >
              <Send size={16} />
              {message.action_card.action.text}
            </button>
          </div>
        )}
        <div className="message-time">
          {message.created_at && format(new Date(message.created_at), 'HH:mm')}
        </div>
      </div>
    </div>
  );
};

export default ChatMessage;

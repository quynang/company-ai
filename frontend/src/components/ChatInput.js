import React, { useState } from 'react';
import { Send, Paperclip } from 'lucide-react';
import './ChatInput.css';

const ChatInput = ({ onSendMessage, disabled = false, placeholder = "Nhập tin nhắn..." }) => {
  const [message, setMessage] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    if (message.trim() && !disabled) {
      onSendMessage(message.trim());
      setMessage('');
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e);
    }
  };

  return (
    <div className="chat-input-container">
      <form onSubmit={handleSubmit} className="chat-input-form">
        <div className="input-wrapper">
          <button
            type="button"
            className="attachment-button"
            disabled={disabled}
            title="Đính kèm file"
          >
            <Paperclip size={20} />
          </button>
          
          <textarea
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder={placeholder}
            disabled={disabled}
            className="message-input"
            rows="1"
          />
          
          <button
            type="submit"
            className="send-button"
            disabled={disabled || !message.trim()}
            title="Gửi tin nhắn"
          >
            <Send size={20} />
          </button>
        </div>
      </form>
    </div>
  );
};

export default ChatInput;

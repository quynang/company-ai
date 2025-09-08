import React, { useState, useEffect, useRef } from "react";
import {
  Menu,
  Bot,
  Settings,
  FileText,
  Tag,
  BookOpen,
  Briefcase,
  Code,
  Database,
  Globe,
  Shield,
  Zap,
  Star,
  Heart,
  Target,
  Lightbulb,
  Rocket,
  Puzzle,
  Gamepad2,
} from "lucide-react";
import ChatMessage from "./ChatMessage";
import ChatInput from "./ChatInput";
import { categoryAPI } from "../services/api";
import "./ChatWindow.css";

const ChatWindow = ({
  session,
  messages,
  onSendMessage,
  isLoading = false,
  onToggleSidebar,
  onSwitchToAdmin,
  setMessages,
  onCreateSessionWithCategory,
}) => {
  const [isTyping, setIsTyping] = useState(false);
  const [categories, setCategories] = useState([]);
  const [categoriesLoading, setCategoriesLoading] = useState(false);

  // Random icons for categories
  const categoryIcons = [
    FileText,
    Tag,
    BookOpen,
    Briefcase,
    Code,
    Database,
    Globe,
    Shield,
    Zap,
    Star,
    Heart,
    Target,
    Lightbulb,
    Rocket,
    Puzzle,
    Gamepad2,
  ];

  // Get random icon for category based on name
  const getCategoryIcon = (categoryName) => {
    const hash = categoryName.split("").reduce((a, b) => {
      a = (a << 5) - a + b.charCodeAt(0);
      return a & a;
    }, 0);
    const IconComponent = categoryIcons[Math.abs(hash) % categoryIcons.length];
    return IconComponent;
  };
  const messagesEndRef = useRef(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages, isTyping]);

  useEffect(() => {
    if (!session) {
      loadCategories();
    }
  }, [session]);

  const loadCategories = async () => {
    try {
      setCategoriesLoading(true);
      const response = await categoryAPI.getCategories();
      setCategories(response.categories || []);
    } catch (error) {
      console.error("Error loading categories:", error);
    } finally {
      setCategoriesLoading(false);
    }
  };

  const handleCategoryClick = async (categoryId) => {
    if (onCreateSessionWithCategory) {
      await onCreateSessionWithCategory(categoryId);
    }
  };

  const handleSendMessage = async (message) => {
    setIsTyping(true);
    try {
      await onSendMessage(message);
    } finally {
      setIsTyping(false);
    }
  };

  const handleWelcomeInput = async (message) => {
    // Create a new session without category when user sends message from welcome screen
    if (!session && message.trim()) {
      setIsTyping(true);
      try {
        // First create a new session, then send the message
        await onSendMessage(message);
      } finally {
        setIsTyping(false);
      }
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
            <p>
              {session ? session.name : "Chọn hoặc tạo cuộc trò chuyện mới"}
            </p>
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
            <div>
              <h2>
                Nhận trợ giúp về <span className="highlight">mọi thứ</span> trên
                Company AI
              </h2>
              <div className="welcome-subtitle">
                <p>Đặt câu hỏi. Tìm câu trả lời.</p>
                <p>Quay lại với hoạt động làm việc.</p>
              </div>
            </div>

            {/* Search input field */}
            <div className="welcome-search-section">
              <div className="search-input-container">
                <svg
                  className="search-icon"
                  width="20"
                  height="20"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <circle cx="11" cy="11" r="8"></circle>
                  <path d="m21 21-4.35-4.35"></path>
                </svg>
                <input
                  type="text"
                  className="search-input"
                  placeholder="Bạn cần trợ giúp về vấn đề gì?"
                  onKeyPress={(e) => {
                    if (e.key === "Enter" && e.target.value.trim()) {
                      handleWelcomeInput(e.target.value.trim());
                      e.target.value = "";
                    }
                  }}
                  disabled={isLoading || isTyping}
                />
              </div>
            </div>

            {/* Suggested questions */}
            {categories.length > 0 && (
              <div className="suggested-questions">
                <h3>Hỏi về chuyên mục nào?</h3>
                <div className="questions-grid">
                  {categories.map((category) => {
                    const IconComponent = getCategoryIcon(category.name);
                    return (
                      <button
                        key={category.id}
                        className="question-card"
                        onClick={() => handleCategoryClick(category.id)}
                        disabled={categoriesLoading}
                      >
                        <span className="question-icon">
                          <IconComponent size={20} />
                        </span>
                        <span className="question-text">{category.name}</span>
                      </button>
                    );
                  })}
                </div>
              </div>
            )}
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
                      role: "assistant",
                      content: result.message,
                      created_at: new Date().toISOString(),
                    };
                    setMessages((prev) => [...prev, successMessage]);
                  } else if (result.error) {
                    // Show error message
                    const errorMessage = {
                      id: Date.now(),
                      role: "assistant",
                      content: result.error,
                      created_at: new Date().toISOString(),
                    };
                    setMessages((prev) => [...prev, errorMessage]);
                  }
                }}
              />
            ))}
            {isTyping && <ChatMessage isLoading={true} />}
            <div ref={messagesEndRef} />
          </div>
        )}
      </div>

      {session && (
        <ChatInput
          onSendMessage={handleSendMessage}
          disabled={isLoading || isTyping}
          placeholder="Nhập câu hỏi của bạn..."
        />
      )}
    </div>
  );
};

export default ChatWindow;

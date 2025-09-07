import React, { useState, useEffect } from 'react';
import SessionList from './components/SessionList';
import ChatWindow from './components/ChatWindow';
import AdminDashboard from './components/AdminDashboard';
import { chatAPI } from './services/api';
import './App.css';

function App() {
  const [sessions, setSessions] = useState([]);
  const [currentSession, setCurrentSession] = useState(null);
  const [messages, setMessages] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isSidebarOpen, setIsSidebarOpen] = useState(true);
  const [error, setError] = useState(null);
  const [currentView, setCurrentView] = useState('chat'); // 'chat' or 'admin'

  // Load sessions on component mount
  useEffect(() => {
    loadSessions();
  }, []);

  // Load messages when session changes
  useEffect(() => {
    if (currentSession) {
      loadMessages(currentSession.id);
    } else {
      setMessages([]);
    }
  }, [currentSession]);

  const loadSessions = async () => {
    try {
      setIsLoading(true);
      const response = await chatAPI.getSessions();
      setSessions(response.sessions || []);
      setError(null);
    } catch (error) {
      console.error('Error loading sessions:', error);
      setError('Không thể tải danh sách cuộc trò chuyện');
    } finally {
      setIsLoading(false);
    }
  };

  const loadMessages = async (sessionId) => {
    try {
      setIsLoading(true);
      const response = await chatAPI.getSession(sessionId);
      setMessages(response.messages || []);
      setError(null);
    } catch (error) {
      console.error('Error loading messages:', error);
      setError('Không thể tải tin nhắn');
    } finally {
      setIsLoading(false);
    }
  };

  const createNewSession = async () => {
    try {
      setIsLoading(true);
      const sessionName = `Cuộc trò chuyện ${new Date().toLocaleString('vi-VN')}`;
      const response = await chatAPI.createSession(sessionName);
      
      if (response.session) {
        const newSession = response.session;
        setSessions(prev => [newSession, ...prev]);
        setCurrentSession(newSession);
        setMessages([]);
        setError(null);
      }
    } catch (error) {
      console.error('Error creating session:', error);
      setError('Không thể tạo cuộc trò chuyện mới');
    } finally {
      setIsLoading(false);
    }
  };

  const selectSession = (session) => {
    setCurrentSession(session);
    if (window.innerWidth <= 768) {
      setIsSidebarOpen(false);
    }
  };

  const deleteSession = async (sessionId) => {
    try {
      setIsLoading(true);
      await chatAPI.deleteSession(sessionId);
      
      // Remove session from list
      setSessions(prev => prev.filter(s => s.id !== sessionId));
      
      // If deleted session was current, clear it
      if (currentSession?.id === sessionId) {
        setCurrentSession(null);
        setMessages([]);
      }
      
      setError(null);
    } catch (error) {
      console.error('Error deleting session:', error);
      setError('Không thể xóa cuộc trò chuyện');
    } finally {
      setIsLoading(false);
    }
  };

  const sendMessage = async (messageText) => {
    if (!currentSession) {
      // Create new session if none exists
      await createNewSession();
      return;
    }

    try {
      // Add user message to UI immediately
      const userMessage = {
        id: Date.now(),
        role: 'user',
        content: messageText,
        created_at: new Date().toISOString(),
      };
      
      setMessages(prev => [...prev, userMessage]);
      
      // Send message to backend
      const response = await chatAPI.sendMessage(currentSession.id, messageText);
      
      // Debug logging
      console.log('Backend response:', response);
      console.log('Response structure:', {
        hasResponse: !!response.response,
        hasMessage: !!(response.response && response.response.message),
        messageContent: response.response?.message?.content,
        messageRole: response.response?.message?.role
      });
      
      if (response.response && response.response.message) {
        // Add assistant response to messages
        const assistantMessage = {
          id: response.response.message.id || Date.now(),
          role: response.response.message.role || 'assistant',
          content: response.response.message.content || '',
          created_at: response.response.message.created_at || new Date().toISOString(),
          action_card: response.response.action_card || null
        };
        
        setMessages(prev => [...prev, assistantMessage]);
        
        // Update session in list (move to top)
        setSessions(prev => {
          const updated = prev.map(s => 
            s.id === currentSession.id 
              ? { ...s, updated_at: new Date().toISOString() }
              : s
          );
          return updated.sort((a, b) => new Date(b.updated_at) - new Date(a.updated_at));
        });
      }
      
      setError(null);
    } catch (error) {
      console.error('Error sending message:', error);
      setError('Không thể gửi tin nhắn');
      
      // Remove the user message from UI if sending failed
      setMessages(prev => prev.filter(m => m.id !== Date.now()));
    }
  };

  const toggleSidebar = () => {
    setIsSidebarOpen(prev => !prev);
  };

  const switchToAdmin = () => {
    setCurrentView('admin');
    setIsSidebarOpen(false);
  };

  const switchToChat = () => {
    setCurrentView('chat');
    setIsSidebarOpen(true);
  };

  return (
    <div className="app">
      {error && (
        <div className="error-banner bg-red-500 text-white p-4 rounded-lg shadow-lg">
          <span>{error}</span>
          <button onClick={() => setError(null)} className="ml-2 hover:bg-red-600 px-2 py-1 rounded">&times;</button>
        </div>
      )}
      
      <div className="app-container">
        <div className={`sidebar ${isSidebarOpen ? 'open' : 'closed'}`}>
          <SessionList
            sessions={sessions}
            currentSession={currentSession}
            onSelectSession={selectSession}
            onNewSession={createNewSession}
            onDeleteSession={deleteSession}
            onSwitchToAdmin={switchToAdmin}
            isLoading={isLoading}
          />
        </div>
        
        <div className="main-content">
          {currentView === 'chat' ? (
            <ChatWindow
              session={currentSession}
              messages={messages}
              onSendMessage={sendMessage}
              isLoading={isLoading}
              onToggleSidebar={toggleSidebar}
              onSwitchToAdmin={switchToAdmin}
              setMessages={setMessages}
            />
          ) : (
            <AdminDashboard />
          )}
        </div>
      </div>
      
      {isSidebarOpen && (
        <div 
          className="sidebar-overlay"
          onClick={() => setIsSidebarOpen(false)}
        />
      )}
    </div>
  );
}

export default App;

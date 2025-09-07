import axios from 'axios';

// Get local IP address for API calls
const getLocalIP = () => {
  // Try to get from environment variable first
  if (import.meta.env.VITE_API_URL) {
    return import.meta.env.VITE_API_URL;
  }
  
  // Get current hostname (IP address when accessed from other devices)
  const currentHost = window.location.hostname;
  
  // If accessed from localhost, use localhost for API
  if (currentHost === 'localhost' || currentHost === '127.0.0.1') {
    return 'http://localhost:8082/api/v1';
  }
  
  // If accessed from IP address, use that IP for API
  return `http://${currentHost}:8082/api/v1`;
};

const API_BASE_URL = getLocalIP();

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Chat API functions
export const chatAPI = {
  // Create a new chat session
  createSession: async (name, userID = null, categoryID = null) => {
    const response = await api.post('/chat/sessions', {
      name,
      user_id: userID,
      category_id: categoryID,
    });
    return response.data;
  },

  // Get all chat sessions
  getSessions: async () => {
    const response = await api.get('/chat/sessions');
    return response.data;
  },

  // Get a specific session with messages
  getSession: async (sessionId) => {
    const response = await api.get(`/chat/sessions/${sessionId}`);
    return response.data;
  },

  // Send a message to a session
  sendMessage: async (sessionId, message) => {
    const response = await api.post(`/chat/sessions/${sessionId}/messages`, {
      message,
    });
    return response.data;
  },

  // Delete a session
  deleteSession: async (sessionId) => {
    const response = await api.delete(`/chat/sessions/${sessionId}`);
    return response.data;
  },
};

// User API functions
export const userAPI = {
  // Get all users
  getUsers: async () => {
    const response = await api.get('/users');
    return response.data;
  },

  // Get user by email
  getUserByEmail: async (email) => {
    const response = await api.get(`/users/by-email?email=${email}`);
    return response.data;
  },

  // Create a new user
  createUser: async (userData) => {
    const response = await api.post('/users', userData);
    return response.data;
  },
};

// Document API functions
export const documentAPI = {
  // Get all documents
  getDocuments: async () => {
    const response = await api.get('/documents');
    return response.data;
  },

  // Get a specific document
  getDocument: async (id) => {
    const response = await api.get(`/documents/${id}`);
    return response.data;
  },

  // Upload a document
  uploadDocument: async (file, categoryIds = []) => {
    const formData = new FormData();
    formData.append('file', file);
    
    // Add category IDs to form data
    categoryIds.forEach(id => {
      formData.append('category_ids', id);
    });
    
    const response = await api.post('/documents/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  // Create document from text content
  uploadDocumentFromText: async (name, content, categoryIds = []) => {
    const response = await api.post('/documents/upload', {
      name: name,
      content: content,
      category_ids: categoryIds
    });
    return response.data;
  },

  // Update document content
  updateDocument: async (id, content) => {
    const response = await api.put(`/documents/${id}`, { content });
    return response.data;
  },

  // Delete a document
  deleteDocument: async (id) => {
    const response = await api.delete(`/documents/${id}`);
    return response.data;
  },

  // Re-embed a document
  reembedDocument: async (id) => {
    const response = await api.post(`/documents/${id}/reembed`);
    return response.data;
  },

  // Re-embed with semantic chunking
  reembedWithSemanticChunking: async (id, config = null) => {
    const response = await api.post('/documents/semantic-reembed', {
      document_id: id,
      config: config
    });
    return response.data;
  },

  // Search documents
  searchDocuments: async (query, limit = 10) => {
    const response = await api.get(`/search?q=${encodeURIComponent(query)}&limit=${limit}`);
    return response.data;
  },

  // Update document categories
  updateDocumentCategories: async (id, categoryIds) => {
    const response = await api.put(`/documents/${id}/categories`, {
      category_ids: categoryIds
    });
    return response.data;
  },

  // Get document categories
  getDocumentCategories: async (id) => {
    const response = await api.get(`/documents/${id}/categories`);
    return response.data;
  },
};

// Category API functions
export const categoryAPI = {
  // Get all categories
  getCategories: async () => {
    const response = await api.get('/categories');
    return response.data;
  },

  // Get a specific category
  getCategory: async (id) => {
    const response = await api.get(`/categories/${id}`);
    return response.data;
  },

  // Create a new category
  createCategory: async (name, description = '') => {
    const response = await api.post('/categories', {
      name: name,
      description: description
    });
    return response.data;
  },

  // Update a category
  updateCategory: async (id, name, description = '') => {
    const response = await api.put(`/categories/${id}`, {
      name: name,
      description: description
    });
    return response.data;
  },

  // Delete a category
  deleteCategory: async (id) => {
    const response = await api.delete(`/categories/${id}`);
    return response.data;
  },

  // Get documents by category
  getDocumentsByCategory: async (id) => {
    const response = await api.get(`/categories/${id}/documents`);
    return response.data;
  },
};

// Health check
export const healthAPI = {
  check: async () => {
    const response = await api.get('/health');
    return response.data;
  },
};

export default api;

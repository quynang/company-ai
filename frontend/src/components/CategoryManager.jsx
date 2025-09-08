import React, { useState, useEffect } from 'react';
import { Search, Plus, Edit, Trash2, FileText, Calendar, Tag, Filter } from 'lucide-react';
import { categoryAPI, documentAPI } from '../services/api';
import './CategoryManager.css';

const CategoryManager = () => {
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [editingCategory, setEditingCategory] = useState(null);
  const [formData, setFormData] = useState({ name: '', description: '' });
  const [searchTerm, setSearchTerm] = useState('');
  const [viewMode, setViewMode] = useState('table'); // Only table view
  const [categoryStats, setCategoryStats] = useState({});

  useEffect(() => {
    loadCategories();
  }, []);

  const loadCategories = async () => {
    try {
      setLoading(true);
      const response = await categoryAPI.getCategories();
      const categoriesData = response.categories || [];
      setCategories(categoriesData);
      
      // Load document counts for each category
      const stats = {};
      for (const category of categoriesData) {
        try {
          const docResponse = await documentAPI.getDocumentsByCategory(category.id);
          stats[category.id] = docResponse.documents?.length || 0;
        } catch (err) {
          console.error(`Error loading documents for category ${category.id}:`, err);
          stats[category.id] = 0;
        }
      }
      setCategoryStats(stats);
      setError(null);
    } catch (err) {
      console.error('Error loading categories:', err);
      setError('Failed to load categories');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateCategory = async (e) => {
    e.preventDefault();
    try {
      await categoryAPI.createCategory(formData.name, formData.description);
      setFormData({ name: '', description: '' });
      setShowCreateForm(false);
      loadCategories();
    } catch (err) {
      console.error('Error creating category:', err);
      setError('Failed to create category');
    }
  };

  const handleUpdateCategory = async (e) => {
    e.preventDefault();
    try {
      await categoryAPI.updateCategory(editingCategory.id, formData.name, formData.description);
      setFormData({ name: '', description: '' });
      setEditingCategory(null);
      loadCategories();
    } catch (err) {
      console.error('Error updating category:', err);
      setError('Failed to update category');
    }
  };

  const handleDeleteCategory = async (categoryId) => {
    if (!window.confirm('Are you sure you want to delete this category? This action cannot be undone.')) {
      return;
    }

    try {
      await categoryAPI.deleteCategory(categoryId);
      loadCategories();
    } catch (err) {
      console.error('Error deleting category:', err);
      setError('Failed to delete category. Make sure no documents or sessions are using this category.');
    }
  };

  const startEdit = (category) => {
    setEditingCategory(category);
    setFormData({ name: category.name, description: category.description || '' });
    setShowCreateForm(false);
  };

  const cancelEdit = () => {
    setEditingCategory(null);
    setFormData({ name: '', description: '' });
    setShowCreateForm(false);
  };

  // Filter categories based on search term
  const filteredCategories = categories.filter(category =>
    category.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    (category.description && category.description.toLowerCase().includes(searchTerm.toLowerCase()))
  );

  // Get category color based on name
  const getCategoryColor = (name) => {
    const colors = [
      'category-badge-blue',
      'category-badge-green',
      'category-badge-purple',
      'category-badge-pink',
      'category-badge-indigo',
      'category-badge-yellow',
      'category-badge-red',
      'category-badge-gray'
    ];
    const hash = name.split('').reduce((a, b) => {
      a = ((a << 5) - a) + b.charCodeAt(0);
      return a & a;
    }, 0);
    return colors[Math.abs(hash) % colors.length];
  };

  if (loading) {
    return (
      <div className="category-manager">
        <div className="category-manager-container">
          <div className="loading-container">
            <div className="loading-skeleton">
              <div className="skeleton-header"></div>
              <div className="skeleton-row"></div>
              <div className="skeleton-row"></div>
              <div className="skeleton-row"></div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="category-manager">
      <div className="category-manager-container">
        {/* Search and Controls */}
        <div className="controls-section">
          <div className="controls-row">
            <div className="search-container">
              <Search className="search-icon" size={20} />
              <input
                type="text"
                placeholder="Search categories..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="search-input"
              />
            </div>
            <button
              onClick={() => setShowCreateForm(true)}
              className="create-btn"
            >
              <Plus size={16} />
              Create Category
            </button>
          </div>
        </div>

        {error && (
          <div className="error-message">
            {error}
          </div>
        )}

        {/* Create/Edit Form */}
        {(showCreateForm || editingCategory) && (
          <div className="form-section">
            <h3 className="form-title">
              {editingCategory ? 'Edit Category' : 'Create New Category'}
            </h3>
            <form onSubmit={editingCategory ? handleUpdateCategory : handleCreateCategory}>
              <div className="form-grid">
                <div className="form-group">
                  <label htmlFor="name" className="form-label">
                    Name *
                  </label>
                  <input
                    type="text"
                    id="name"
                    required
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    className="form-input"
                    placeholder="Category name"
                  />
                </div>
                <div className="form-group">
                  <label htmlFor="description" className="form-label">
                    Description
                  </label>
                  <input
                    type="text"
                    id="description"
                    value={formData.description}
                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                    className="form-input"
                    placeholder="Category description"
                  />
                </div>
              </div>
              <div className="form-buttons">
                <button
                  type="submit"
                  className="btn btn-primary"
                >
                  {editingCategory ? 'Update' : 'Create'}
                </button>
                <button
                  type="button"
                  onClick={cancelEdit}
                  className="btn btn-secondary"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        )}

        {/* Categories Table */}
        <div className="table-container">
          {filteredCategories.length === 0 ? (
            <div className="empty-state">
              <div className="empty-state-text">
                {searchTerm ? 'No categories found matching your search.' : 'No categories found. Create your first category to get started.'}
              </div>
            </div>
          ) : (
            <table className="categories-table">
              <thead className="table-header">
                <tr>
                  <th>Category</th>
                  <th>Description</th>
                  <th>Documents</th>
                  <th>Created</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {filteredCategories.map((category) => (
                  <tr key={category.id} className="table-row">
                    <td className="table-cell">
                      <div className={`category-badge ${getCategoryColor(category.name)}`}>
                        <Tag size={12} />
                        {category.name}
                      </div>
                    </td>
                    <td className="table-cell">
                      <div className="description-cell">
                        {category.description || 'No description'}
                      </div>
                    </td>
                    <td className="table-cell">
                      <div className="stats-cell">
                        <FileText size={14} />
                        {categoryStats[category.id] || 0}
                      </div>
                    </td>
                    <td className="table-cell">
                      <div className="date-cell">
                        {new Date(category.created_at).toLocaleDateString()}
                      </div>
                    </td>
                    <td className="table-cell">
                      <div className="actions-cell">
                        <button
                          onClick={() => startEdit(category)}
                          className="action-btn action-btn-edit"
                        >
                          <Edit size={14} />
                          Edit
                        </button>
                        <button
                          onClick={() => handleDeleteCategory(category.id)}
                          className="action-btn action-btn-delete"
                        >
                          <Trash2 size={14} />
                          Delete
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  );
};

export default CategoryManager;

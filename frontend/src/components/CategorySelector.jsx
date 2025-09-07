import React, { useState, useEffect } from 'react';
import { categoryAPI } from '../services/api';
import './CategorySelector.css';

const CategorySelector = ({ selectedCategoryId, onCategoryChange, className = '' }) => {
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    loadCategories();
  }, []);

  const loadCategories = async () => {
    try {
      setLoading(true);
      const response = await categoryAPI.getCategories();
      setCategories(response.categories || []);
      setError(null);
    } catch (err) {
      console.error('Error loading categories:', err);
      setError('Failed to load categories');
    } finally {
      setLoading(false);
    }
  };

  const handleCategoryChange = (e) => {
    const categoryId = e.target.value === '' ? null : e.target.value;
    onCategoryChange(categoryId);
  };

  if (loading) {
    return (
      <div className={`${className}`}>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Category (Optional)
        </label>
        <div className="animate-pulse bg-gray-200 h-10 rounded-md"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={`${className}`}>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Category (Optional)
        </label>
        <div className="text-red-600 text-sm">{error}</div>
      </div>
    );
  }

  return (
    <div className={`${className}`}>
      <label htmlFor="category-select" className="block text-sm font-medium text-gray-700 mb-2">
        Category (Optional)
      </label>
      <select
        id="category-select"
        value={selectedCategoryId || ''}
        onChange={handleCategoryChange}
        className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
      >
        <option value="">All Categories</option>
        {categories.map((category) => (
          <option key={category.id} value={category.id}>
            {category.name}
          </option>
        ))}
      </select>
      {selectedCategoryId && (
        <p className="mt-1 text-sm text-gray-500">
          Chat will only search documents in the selected category
        </p>
      )}
    </div>
  );
};

export default CategorySelector;

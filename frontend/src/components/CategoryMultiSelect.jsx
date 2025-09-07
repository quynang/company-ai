import React, { useState, useEffect } from 'react';
import { categoryAPI } from '../services/api';
import './CategoryMultiSelect.css';

const CategoryMultiSelect = ({ selectedCategoryIds = [], onCategoryChange, className = '' }) => {
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isOpen, setIsOpen] = useState(false);

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

  const handleCategoryToggle = (categoryId) => {
    const newSelectedIds = selectedCategoryIds.includes(categoryId)
      ? selectedCategoryIds.filter(id => id !== categoryId)
      : [...selectedCategoryIds, categoryId];
    
    onCategoryChange(newSelectedIds);
  };

  const handleSelectAll = () => {
    if (selectedCategoryIds.length === categories.length) {
      onCategoryChange([]);
    } else {
      onCategoryChange(categories.map(cat => cat.id));
    }
  };

  const getSelectedCategories = () => {
    return categories.filter(cat => selectedCategoryIds.includes(cat.id));
  };

  if (loading) {
    return (
      <div className={`category-multi-select ${className}`}>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Categories
        </label>
        <div className="animate-pulse bg-gray-200 h-10 rounded-md"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={`category-multi-select ${className}`}>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Categories
        </label>
        <div className="text-red-600 text-sm">{error}</div>
      </div>
    );
  }

  return (
    <div className={`category-multi-select ${className}`}>
      <label className="block text-sm font-medium text-gray-700 mb-2">
        Categories
      </label>
      
      <div className="relative">
        <button
          type="button"
          onClick={() => setIsOpen(!isOpen)}
          className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-left"
        >
          {selectedCategoryIds.length === 0 ? (
            <span className="text-gray-500">Select categories...</span>
          ) : (
            <span>
              {selectedCategoryIds.length} categor{selectedCategoryIds.length === 1 ? 'y' : 'ies'} selected
            </span>
          )}
          <span className="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
            <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
            </svg>
          </span>
        </button>

        {isOpen && (
          <div className="absolute z-10 mt-1 w-full bg-white shadow-lg max-h-60 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none">
            <div className="px-3 py-2 border-b border-gray-200">
              <button
                type="button"
                onClick={handleSelectAll}
                className="text-sm text-blue-600 hover:text-blue-800"
              >
                {selectedCategoryIds.length === categories.length ? 'Deselect All' : 'Select All'}
              </button>
            </div>
            
            {categories.map((category) => (
              <div
                key={category.id}
                className="px-3 py-2 hover:bg-gray-100 cursor-pointer"
                onClick={() => handleCategoryToggle(category.id)}
              >
                <div className="flex items-center">
                  <input
                    type="checkbox"
                    checked={selectedCategoryIds.includes(category.id)}
                    onChange={() => handleCategoryToggle(category.id)}
                    className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded mr-3"
                  />
                  <div className="flex-1">
                    <div className="text-sm font-medium text-gray-900">{category.name}</div>
                    {category.description && (
                      <div className="text-xs text-gray-500">{category.description}</div>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {selectedCategoryIds.length > 0 && (
        <div className="mt-2">
          <div className="text-sm text-gray-600 mb-1">Selected categories:</div>
          <div className="flex flex-wrap gap-1">
            {getSelectedCategories().map((category) => (
              <span
                key={category.id}
                className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
              >
                {category.name}
                <button
                  type="button"
                  onClick={() => handleCategoryToggle(category.id)}
                  className="ml-1 inline-flex items-center justify-center w-4 h-4 rounded-full text-blue-400 hover:bg-blue-200 hover:text-blue-500"
                >
                  <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </span>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default CategoryMultiSelect;

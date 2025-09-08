import React, { useState, useEffect } from 'react';
import Select from 'react-select';
import { categoryAPI } from '../services/api';
import './CategoryMultiSelect.css';

const CategoryMultiSelect = ({ selectedCategoryIds = [], onCategoryChange, className = '' }) => {
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

  // Get category color based on name
  const getCategoryColor = (name) => {
    const colors = [
      'category-option-badge-blue',
      'category-option-badge-green',
      'category-option-badge-purple',
      'category-option-badge-pink',
      'category-option-badge-indigo',
      'category-option-badge-yellow',
      'category-option-badge-red',
      'category-option-badge-gray'
    ];
    const hash = name.split('').reduce((a, b) => {
      a = ((a << 5) - a) + b.charCodeAt(0);
      return a & a;
    }, 0);
    return colors[Math.abs(hash) % colors.length];
  };

  // Convert categories to react-select format
  const selectOptions = categories.map(category => ({
    value: category.id,
    label: category.name,
    description: category.description,
    color: getCategoryColor(category.name)
  }));

  // Get selected options
  const selectedOptions = selectOptions.filter(option => 
    selectedCategoryIds.includes(option.value)
  );

  // Handle selection change
  const handleChange = (selectedOptions) => {
    const newSelectedIds = selectedOptions ? selectedOptions.map(option => option.value) : [];
    onCategoryChange(newSelectedIds);
  };

  // Custom option component
  const CustomOption = ({ innerProps, innerRef, data, isFocused, isSelected }) => (
    <div
      ref={innerRef}
      {...innerProps}
      className={`category-option-content ${isFocused ? 'focused' : ''} ${isSelected ? 'selected' : ''}`}
    >
      <div className={`category-option-badge ${data.color}`}></div>
      <div className="category-option-text">
        <div className="category-option-name">{data.label}</div>
      </div>
    </div>
  );

  // Custom placeholder component
  const CustomPlaceholder = () => (
    <div style={{ color: '#9ca3af', fontWeight: 500 }}>
      Select categories...
    </div>
  );

  // Custom loading message
  const LoadingMessage = () => (
    <div className="category-select-loading">
      Loading categories...
    </div>
  );

  // Custom no options message
  const NoOptionsMessage = () => (
    <div style={{ padding: '20px', textAlign: 'center', color: '#6b7280', fontStyle: 'italic' }}>
      No categories found
    </div>
  );

  if (loading) {
    return (
      <div className={`category-multi-select ${className}`}>
        <label className="category-multi-select-label">
          Categories
        </label>
        <div className="category-select-loading">
          Loading categories...
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={`category-multi-select ${className}`}>
        <label className="category-multi-select-label">
          Categories
        </label>
        <div className="category-select-error">
          {error}
        </div>
      </div>
    );
  }

  return (
    <div className={`category-multi-select ${className}`}>
      <label className="category-multi-select-label">
        Categories
      </label>
      
      <div className="category-select-container">
        <Select
          isMulti
          options={selectOptions}
          value={selectedOptions}
          onChange={handleChange}
          placeholder="Select categories..."
          className="react-select-container"
          classNamePrefix="react-select"
          menuPortalTarget={document.body}
          menuPosition="fixed"
          components={{
            Option: CustomOption,
            LoadingMessage,
            NoOptionsMessage
          }}
          styles={{
            control: (provided, state) => ({
              ...provided,
              minHeight: '48px',
              border: '2px solid #e5e7eb',
              borderRadius: '12px',
              boxShadow: state.isFocused 
                ? '0 0 0 3px rgba(102, 126, 234, 0.1), 0 4px 12px rgba(102, 126, 234, 0.15)' 
                : '0 2px 4px rgba(0, 0, 0, 0.05)',
              '&:hover': {
                borderColor: '#667eea',
                boxShadow: '0 4px 12px rgba(102, 126, 234, 0.15)',
                transform: 'translateY(-1px)'
              }
            }),
            multiValue: (provided) => ({
              ...provided,
              background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
              borderRadius: '20px',
              border: 'none',
              boxShadow: '0 2px 4px rgba(102, 126, 234, 0.2)'
            }),
            multiValueLabel: (provided) => ({
              ...provided,
              color: 'white',
              fontWeight: '600',
              fontSize: '0.75rem',
              textTransform: 'uppercase',
              letterSpacing: '0.025em'
            }),
            multiValueRemove: (provided) => ({
              ...provided,
              color: 'rgba(255, 255, 255, 0.8)',
              borderRadius: '50%',
              '&:hover': {
                backgroundColor: 'rgba(255, 255, 255, 0.2)',
                color: 'white',
                transform: 'scale(1.1)'
              }
            }),
            menu: (provided) => ({
              ...provided,
              background: 'white',
              border: '2px solid #e5e7eb',
              borderRadius: '12px',
              boxShadow: '0 20px 40px rgba(0, 0, 0, 0.1)',
              backdropFilter: 'blur(10px)',
              marginTop: '8px',
              overflow: 'hidden',
              zIndex: 9999
            }),
            menuList: (provided) => ({
              ...provided,
              padding: '12px'
            }),
            menuPortal: (provided) => ({
              ...provided,
              zIndex: 9999
            }),
            option: (provided, state) => ({
              ...provided,
              padding: '12px 16px',
              borderRadius: '10px',
              margin: '4px 0',
              fontSize: '0.875rem',
              fontWeight: '500',
              color: state.isFocused ? 'white' : '#374151',
              background: state.isFocused 
                ? 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
                : state.isSelected 
                  ? 'linear-gradient(135deg, #10b981 0%, #059669 100%)'
                  : 'transparent',
              minHeight: '44px',
              display: 'flex',
              alignItems: 'center',
              border: '1px solid transparent',
              cursor: 'pointer',
              transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
              '&:hover': {
                background: state.isFocused 
                  ? 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
                  : 'linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%)',
                transform: 'translateX(6px) scale(1.02)',
                boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)',
                borderColor: 'rgba(102, 126, 234, 0.2)'
              }
            })
          }}
          theme={(theme) => ({
            ...theme,
            colors: {
              ...theme.colors,
              primary: '#667eea',
              primary75: '#8b5cf6',
              primary50: '#a5b4fc',
              primary25: '#c7d2fe',
              danger: '#ef4444',
              dangerLight: '#fecaca',
              neutral0: '#ffffff',
              neutral5: '#f9fafb',
              neutral10: '#f3f4f6',
              neutral20: '#e5e7eb',
              neutral30: '#d1d5db',
              neutral40: '#9ca3af',
              neutral50: '#6b7280',
              neutral60: '#4b5563',
              neutral70: '#374151',
              neutral80: '#1f2937',
              neutral90: '#111827'
            }
          })}
        />
      </div>
    </div>
  );
};

export default CategoryMultiSelect;


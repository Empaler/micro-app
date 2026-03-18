import { useState } from 'react';
import './BookForm.css';

function BookForm({ book, onSubmit, onClose }) {
  const [formData, setFormData] = useState({
    title: book?.title || '',
    author: book?.author || '',
    releaseYear: book?.releaseYear || new Date().getFullYear(),
    rating: book?.rating || 5,
  });

  const [errors, setErrors] = useState({});

  const validate = () => {
    const newErrors = {};
    if (!formData.title.trim()) newErrors.title = 'Title is required';
    if (!formData.author.trim()) newErrors.author = 'Author is required';
    if (formData.releaseYear < 1000 || formData.releaseYear > new Date().getFullYear()) {
      newErrors.releaseYear = 'Year must be between 1000 and current year';
    }
    if (formData.rating < 0 || formData.rating > 10) {
      newErrors.rating = 'Rating must be between 0 and 10';
    }
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!validate()) return;
    onSubmit(formData);
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: name === 'rating' || name === 'releaseYear' ? parseFloat(value) || 0 : value,
    }));
  };

  return (
    <div className="modal-overlay">
      <div className="modal">
        <div className="modal-header">
          <h2>{book ? 'Edit Book' : 'Add New Book'}</h2>
          <button className="close-btn" onClick={onClose}>×</button>
        </div>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Title *</label>
            <input
              type="text"
              name="title"
              value={formData.title}
              onChange={handleChange}
              className={errors.title ? 'error' : ''}
            />
            {errors.title && <span className="error-text">{errors.title}</span>}
          </div>

          <div className="form-group">
            <label>Author *</label>
            <input
              type="text"
              name="author"
              value={formData.author}
              onChange={handleChange}
              className={errors.author ? 'error' : ''}
            />
            {errors.author && <span className="error-text">{errors.author}</span>}
          </div>

          <div className="form-row">
            <div className="form-group">
              <label>Release Year *</label>
              <input
                type="number"
                name="releaseYear"
                value={formData.releaseYear}
                onChange={handleChange}
                className={errors.releaseYear ? 'error' : ''}
              />
              {errors.releaseYear && <span className="error-text">{errors.releaseYear}</span>}
            </div>

            <div className="form-group">
              <label>Rating * (0-10)</label>
              <input
                type="number"
                name="rating"
                step="0.1"
                min="0"
                max="10"
                value={formData.rating}
                onChange={handleChange}
                className={errors.rating ? 'error' : ''}
              />
              {errors.rating && <span className="error-text">{errors.rating}</span>}
            </div>
          </div>

          <div className="form-actions">
            <button type="button" className="btn btn-secondary" onClick={onClose}>
              Cancel
            </button>
            <button type="submit" className="btn btn-primary">
              {book ? 'Update' : 'Add'} Book
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default BookForm;

import { useState } from 'react';
import './MovieForm.css';

function MovieForm({ movie, onSubmit, onClose }) {
  const [formData, setFormData] = useState({
    title: movie?.title || '',
    year: movie?.year || new Date().getFullYear(),
    type: movie?.type || 'movie',
    resolution: movie?.resolution || 'FHD',
    actors: movie?.actors || '',
    rating: movie?.rating || 5,
    isAdult: movie?.isAdult || false,
  });

  const [errors, setErrors] = useState({});

  const validate = () => {
    const newErrors = {};
    if (!formData.title.trim()) newErrors.title = 'Title is required';
    if (formData.year < 1888 || formData.year > new Date().getFullYear()) {
      newErrors.year = 'Year must be between 1888 and current year';
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
    const { name, value, type, checked } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value,
    }));
  };

  return (
    <div className="modal-overlay">
      <div className="modal">
        <div className="modal-header">
          <h2>{movie ? 'Edit Movie' : 'Add New Movie'}</h2>
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

          <div className="form-row">
            <div className="form-group">
              <label>Year *</label>
              <input
                type="number"
                name="year"
                value={formData.year}
                onChange={handleChange}
                className={errors.year ? 'error' : ''}
              />
              {errors.year && <span className="error-text">{errors.year}</span>}
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

          <div className="form-row">
            <div className="form-group">
              <label>Type *</label>
              <select name="type" value={formData.type} onChange={handleChange}>
                <option value="movie">Movie</option>
                <option value="series">Series</option>
              </select>
            </div>

            <div className="form-group">
              <label>Resolution *</label>
              <select name="resolution" value={formData.resolution} onChange={handleChange}>
                <option value="SD">SD</option>
                <option value="HD">HD</option>
                <option value="FHD">FHD</option>
                <option value="4K">4K</option>
              </select>
            </div>
          </div>

          <div className="form-group">
            <label>Actors</label>
            <input
              type="text"
              name="actors"
              value={formData.actors}
              onChange={handleChange}
              placeholder="Comma-separated actor names"
            />
          </div>

          <div className="form-actions">
            <button type="button" className="btn btn-secondary" onClick={onClose}>
              Cancel
            </button>
            <button type="submit" className="btn btn-primary">
              {movie ? 'Update' : 'Add'} Movie
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default MovieForm;

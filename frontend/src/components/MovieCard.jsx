import './MovieCard.css';

function MovieCard({ movie, onView, onEdit, onDelete }) {
  return (
    <div className="movie-card">
      <div className="movie-card-header">
        <h3>{movie.title}</h3>
        <span className="movie-rating">★ {movie.rating}</span>
      </div>
      <div className="movie-card-body">
        <p><strong>Year:</strong> {movie.year}</p>
        <p><strong>Type:</strong> {movie.type}</p>
        <p><strong>Resolution:</strong> {movie.resolution}</p>
        {movie.actors && <p><strong>Actors:</strong> {movie.actors}</p>}
      </div>
      <div className="movie-card-actions">
        <button className="btn btn-primary" onClick={onView}>View Details</button>
        <button className="btn btn-secondary" onClick={onEdit}>Edit</button>
        <button className="btn btn-danger" onClick={onDelete}>Delete</button>
      </div>
    </div>
  );
}

export default MovieCard;

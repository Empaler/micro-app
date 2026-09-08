import './BookCard.css';

function BookCard({ book, onView, onEdit, onDelete }) {
  return (
    <div className="book-card">
      <div className="book-card-header">
        <h3>{book.title}</h3>
        <span className="book-rating">★ {book.rating}</span>
      </div>
      <div className="book-card-body">
        <p><strong>Author:</strong> {book.author}</p>
        <p><strong>Release Year:</strong> {book.releaseYear}</p>
      </div>
      <div className="book-card-actions">
        <button className="btn btn-primary" onClick={onView}>View Details</button>
        <button className="btn btn-secondary" onClick={onEdit}>Edit</button>
        <button className="btn btn-danger" onClick={onDelete}>Delete</button>
      </div>
    </div>
  );
}

export default BookCard;

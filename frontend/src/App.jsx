import { useState, useEffect } from 'react';
import { movieApi, bookApi } from './api';
import MovieForm from './components/MovieForm';
import MovieCard from './components/MovieCard';
import BookForm from './components/BookForm';
import BookCard from './components/BookCard';
import './App.css';

function App() {
  const [activeTab, setActiveTab] = useState('movies');
  const [movies, setMovies] = useState([]);
  const [books, setBooks] = useState([]);
  const [popularMovies, setPopularMovies] = useState([]);
  const [popularBooks, setPopularBooks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showMovieForm, setShowMovieForm] = useState(false);
  const [showBookForm, setShowBookForm] = useState(false);
  const [editingMovie, setEditingMovie] = useState(null);
  const [editingBook, setEditingBook] = useState(null);
  const [selectedDetail, setSelectedDetail] = useState(null);

  useEffect(() => {
    fetchData();
  }, [activeTab]);

  const fetchData = async () => {
    setLoading(true);
    try {
      if (activeTab === 'movies') {
        const [moviesResponse, popularResponse] = await Promise.all([
          movieApi.getAll(),
          movieApi.getMostLookedUp(),
        ]);
        setMovies(moviesResponse.data.data || []);
        setPopularMovies(popularResponse.data.data || []);
      } else {
        const [booksResponse, popularResponse] = await Promise.all([
          bookApi.getAll(),
          bookApi.getMostLookedUp(),
        ]);
        setBooks(booksResponse.data.data || []);
        setPopularBooks(popularResponse.data.data || []);
      }
      setError(null);
    } catch (err) {
      setError(`Failed to fetch ${activeTab}`);
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleAddMovie = async (movieData) => {
    try {
      await movieApi.create(movieData);
      setShowMovieForm(false);
      fetchData();
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to add movie');
    }
  };

  const handleDeleteMovie = async (id) => {
    if (!window.confirm('Are you sure you want to delete this movie?')) return;
    try {
      await movieApi.delete(id);
      fetchData();
    } catch (err) {
      alert('Failed to delete movie');
    }
  };

  const handleEditMovie = (movie) => {
    setEditingMovie(movie);
    setShowMovieForm(true);
  };

  const handleUpdateMovie = async (movieData) => {
    try {
      await movieApi.update(editingMovie.id, movieData);
      setShowMovieForm(false);
      setEditingMovie(null);
      fetchData();
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to update movie');
    }
  };

  const handleAddBook = async (bookData) => {
    try {
      await bookApi.create(bookData);
      setShowBookForm(false);
      fetchData();
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to add book');
    }
  };

  const handleDeleteBook = async (id) => {
    if (!window.confirm('Are you sure you want to delete this book?')) return;
    try {
      await bookApi.delete(id);
      fetchData();
    } catch (err) {
      alert('Failed to delete book');
    }
  };

  const handleEditBook = (book) => {
    setEditingBook(book);
    setShowBookForm(true);
  };

  const handleUpdateBook = async (bookData) => {
    try {
      await bookApi.update(editingBook.id, bookData);
      setShowBookForm(false);
      setEditingBook(null);
      fetchData();
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to update book');
    }
  };

  const handleCloseMovieForm = () => {
    setShowMovieForm(false);
    setEditingMovie(null);
  };

  const handleCloseBookForm = () => {
    setShowBookForm(false);
    setEditingBook(null);
  };

  const handleViewMovie = async (id) => {
    try {
      const response = await movieApi.getById(id);
      setSelectedDetail({ type: 'Movie', item: response.data.data });
      const popularResponse = await movieApi.getMostLookedUp();
      setPopularMovies(popularResponse.data.data || []);
    } catch (err) {
      setError('Failed to fetch movie details');
    }
  };

  const handleViewBook = async (id) => {
    try {
      const response = await bookApi.getById(id);
      setSelectedDetail({ type: 'Book', item: response.data.data });
      const popularResponse = await bookApi.getMostLookedUp();
      setPopularBooks(popularResponse.data.data || []);
    } catch (err) {
      setError('Failed to fetch book details');
    }
  };

  const renderPopularItems = () => {
    const items = activeTab === 'movies' ? popularMovies : popularBooks;
    const allItems = activeTab === 'movies' ? movies : books;

    return (
      <div className="popular-panel">
        <h2>Most looked up</h2>

        {!items || items.length === 0 ? (
          <div className="popular-empty">
            No popular items yet. Open a {activeTab === 'movies' ? 'movie' : 'book'} detail to start ranking.
          </div>
        ) : (
          <div className="popular-list">
            {items.map((item, index) => {
              const match = allItems.find((entry) => Number(entry.id) === Number(item.id));
              const title = match ? match.title || match.name || 'Unknown title' : `Item #${item.id}`;

              return (
                <div key={`${activeTab}-${item.id}`} className="popular-item">
                  <span className="popular-rank">#{index + 1}</span>
                  <span className="popular-title">{title}</span>
                  <span className="popular-score">{item.score}</span>
                </div>
              );
            })}
          </div>
        )}
      </div>
    );
  };

  return (
    <div className="app">
      <header className="header">
        <h1>Movie Database</h1>
      </header>

      <div className="tabs">
        <button
          className={`tab ${activeTab === 'movies' ? 'active' : ''}`}
          onClick={() => setActiveTab('movies')}
        >
          Movies
        </button>
        <button
          className={`tab ${activeTab === 'books' ? 'active' : ''}`}
          onClick={() => setActiveTab('books')}
        >
          Books
        </button>
      </div>

      <div className="tab-content">
        {activeTab === 'movies' && (
          <>
            {renderPopularItems()}

            <div className="tab-header">
              <button className="btn btn-primary" onClick={() => setShowMovieForm(true)}>
                Add Movie
              </button>
            </div>

            {error && <div className="error">{error}</div>}

            {loading ? (
              <div className="loading">Loading...</div>
            ) : movies.length === 0 ? (
              <div className="empty">No movies yet. Add your first movie!</div>
            ) : (
              <div className="grid">
                {movies.map((movie) => (
                  <MovieCard
                    key={movie.id}
                    movie={movie}
                    onView={() => handleViewMovie(movie.id)}
                    onEdit={() => handleEditMovie(movie)}
                    onDelete={() => handleDeleteMovie(movie.id)}
                  />
                ))}
              </div>
            )}

            {showMovieForm && (
              <MovieForm
                movie={editingMovie}
                onSubmit={editingMovie ? handleUpdateMovie : handleAddMovie}
                onClose={handleCloseMovieForm}
              />
            )}
          </>
        )}

        {activeTab === 'books' && (
          <>
            {renderPopularItems()}

            <div className="tab-header">
              <button className="btn btn-primary" onClick={() => setShowBookForm(true)}>
                Add Book
              </button>
            </div>

            {error && <div className="error">{error}</div>}

            {loading ? (
              <div className="loading">Loading...</div>
            ) : books.length === 0 ? (
              <div className="empty">No books yet. Add your first book!</div>
            ) : (
              <div className="grid">
                {books.map((book) => (
                  <BookCard
                    key={book.id}
                    book={book}
                    onView={() => handleViewBook(book.id)}
                    onEdit={() => handleEditBook(book)}
                    onDelete={() => handleDeleteBook(book.id)}
                  />
                ))}
              </div>
            )}

            {showBookForm && (
              <BookForm
                book={editingBook}
                onSubmit={editingBook ? handleUpdateBook : handleAddBook}
                onClose={handleCloseBookForm}
              />
            )}
          </>
        )}
      </div>

      {selectedDetail && (
        <div className="detail-modal" role="dialog" aria-modal="true" aria-label={`${selectedDetail.type} details`}>
          <div className="detail-modal-content">
            <button className="detail-modal-close" onClick={() => setSelectedDetail(null)} aria-label="Close details">
              ×
            </button>
            <h2>{selectedDetail.item.title}</h2>
            <p><strong>Type:</strong> {selectedDetail.type}</p>
            {selectedDetail.type === 'Movie' ? (
              <>
                <p><strong>Year:</strong> {selectedDetail.item.year}</p>
                <p><strong>Resolution:</strong> {selectedDetail.item.resolution}</p>
                {selectedDetail.item.actors && <p><strong>Actors:</strong> {selectedDetail.item.actors}</p>}
              </>
            ) : (
              <>
                <p><strong>Author:</strong> {selectedDetail.item.author}</p>
                <p><strong>Release Year:</strong> {selectedDetail.item.releaseYear}</p>
              </>
            )}
          </div>
        </div>
      )}
    </div>
  );
}

export default App;

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
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showMovieForm, setShowMovieForm] = useState(false);
  const [showBookForm, setShowBookForm] = useState(false);
  const [editingMovie, setEditingMovie] = useState(null);
  const [editingBook, setEditingBook] = useState(null);

  useEffect(() => {
    fetchData();
  }, [activeTab]);

  const fetchData = async () => {
    setLoading(true);
    try {
      if (activeTab === 'movies') {
        const response = await movieApi.getAll();
        setMovies(response.data.data || []);
      } else {
        const response = await bookApi.getAll();
        setBooks(response.data.data || []);
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
    </div>
  );
}

export default App;

package movies

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockMovieRepo struct {
	movies map[int64]*Movie
	nextID int64
}

func newMockMovieRepo() *mockMovieRepo {
	return &mockMovieRepo{movies: make(map[int64]*Movie), nextID: 1}
}

func (m *mockMovieRepo) Create(ctx context.Context, movie *Movie) error {
	if movie.ID == 0 {
		movie.ID = m.nextID
		m.nextID++
	}
	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()
	m.movies[movie.ID] = movie
	return nil
}

func (m *mockMovieRepo) GetByID(ctx context.Context, id int64) (*Movie, error) {
	if movie, ok := m.movies[id]; ok {
		return movie, nil
	}
	return nil, ErrMovieNotFound
}

func (m *mockMovieRepo) GetAll(ctx context.Context) ([]Movie, error) {
	result := make([]Movie, 0, len(m.movies))
	for _, movie := range m.movies {
		result = append(result, *movie)
	}
	return result, nil
}

func (m *mockMovieRepo) Update(ctx context.Context, movie *Movie) error {
	if _, ok := m.movies[movie.ID]; !ok {
		return ErrMovieNotFound
	}
	m.movies[movie.ID] = movie
	return nil
}

func (m *mockMovieRepo) Delete(ctx context.Context, id int64) error {
	if _, ok := m.movies[id]; !ok {
		return ErrMovieNotFound
	}
	delete(m.movies, id)
	return nil
}

func TestMovieService_CreateMovie(t *testing.T) {
	tests := []struct {
		name    string
		movie   Movie
		wantErr bool
	}{
		{
			name: "valid movie",
			movie: Movie{
				Title:      "The Matrix",
				Year:       1999,
				Type:       MovieTypeMovie,
				Resolution: ResolutionFHD,
				Rating:     8.7,
			},
			wantErr: false,
		},
		{
			name: "missing title",
			movie: Movie{
				Title:      "",
				Year:       1999,
				Type:       MovieTypeMovie,
				Resolution: ResolutionFHD,
				Rating:     8.7,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockMovieRepo()
			svc := NewMovieService(repo)

			err := svc.CreateMovie(context.Background(), &tt.movie)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateMovie() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMovieService_GetMovie(t *testing.T) {
	ctx := context.Background()
	repo := newMockMovieRepo()
	svc := NewMovieService(repo)

	t.Run("movie not found", func(t *testing.T) {
		_, err := svc.GetMovie(ctx, 999)
		if !errors.Is(err, ErrMovieNotFound) {
			t.Errorf("GetMovie() error = %v, want %v", err, ErrMovieNotFound)
		}
	})

	t.Run("get existing movie", func(t *testing.T) {
		movie := &Movie{
			Title:      "Test",
			Year:       2000,
			Type:       MovieTypeMovie,
			Resolution: ResolutionHD,
			Rating:     7.0,
		}
		repo.Create(ctx, movie)

		got, err := svc.GetMovie(ctx, movie.ID)
		if err != nil {
			t.Fatalf("GetMovie() error = %v", err)
		}
		if got.Title != movie.Title {
			t.Errorf("GetMovie() title = %v, want %v", got.Title, movie.Title)
		}
	})
}

func TestMovieService_GetAllMovies(t *testing.T) {
	ctx := context.Background()
	repo := newMockMovieRepo()
	svc := NewMovieService(repo)

	repo.Create(ctx, &Movie{
		Title:      "Movie 1",
		Year:       2000,
		Type:       MovieTypeMovie,
		Resolution: ResolutionHD,
		Rating:     7.0,
	})
	repo.Create(ctx, &Movie{
		Title:      "Movie 2",
		Year:       2001,
		Type:       MovieTypeSeries,
		Resolution: ResolutionFHD,
		Rating:     8.0,
	})

	movies, err := svc.GetAllMovies(ctx)
	if err != nil {
		t.Fatalf("GetAllMovies() error = %v", err)
	}
	if len(movies) != 2 {
		t.Errorf("GetAllMovies() got %d movies, want 2", len(movies))
	}
}

func TestMovieService_UpdateMovie(t *testing.T) {
	ctx := context.Background()
	repo := newMockMovieRepo()
	svc := NewMovieService(repo)

	t.Run("update non-existent movie", func(t *testing.T) {
		err := svc.UpdateMovie(ctx, &Movie{
			ID:         999,
			Title:      "Updated",
			Year:       2000,
			Type:       MovieTypeMovie,
			Resolution: ResolutionHD,
			Rating:     7.0,
		})
		if !errors.Is(err, ErrMovieNotFound) {
			t.Errorf("UpdateMovie() error = %v, want %v", err, ErrMovieNotFound)
		}
	})

	t.Run("update existing movie", func(t *testing.T) {
		movie := &Movie{
			Title:      "Original",
			Year:       2000,
			Type:       MovieTypeMovie,
			Resolution: ResolutionHD,
			Rating:     7.0,
		}
		repo.Create(ctx, movie)

		movie.Title = "Updated Title"
		err := svc.UpdateMovie(ctx, movie)
		if err != nil {
			t.Fatalf("UpdateMovie() error = %v", err)
		}

		got, _ := repo.GetByID(ctx, movie.ID)
		if got.Title != "Updated Title" {
			t.Errorf("UpdateMovie() title = %v, want Updated Title", got.Title)
		}
	})

	t.Run("update with invalid data", func(t *testing.T) {
		movie := &Movie{
			Title:      "Test",
			Year:       2000,
			Type:       MovieTypeMovie,
			Resolution: ResolutionHD,
			Rating:     7.0,
		}
		repo.Create(ctx, movie)

		movie.Title = ""
		err := svc.UpdateMovie(ctx, movie)
		if err != ErrTitleRequired {
			t.Errorf("UpdateMovie() error = %v, want %v", err, ErrTitleRequired)
		}
	})
}

func TestMovieService_DeleteMovie(t *testing.T) {
	ctx := context.Background()
	repo := newMockMovieRepo()
	svc := NewMovieService(repo)

	t.Run("delete non-existent movie", func(t *testing.T) {
		err := svc.DeleteMovie(ctx, 999)
		if !errors.Is(err, ErrMovieNotFound) {
			t.Errorf("DeleteMovie() error = %v, want %v", err, ErrMovieNotFound)
		}
	})

	t.Run("delete existing movie", func(t *testing.T) {
		movie := &Movie{
			Title:      "Test",
			Year:       2000,
			Type:       MovieTypeMovie,
			Resolution: ResolutionHD,
			Rating:     7.0,
		}
		repo.Create(ctx, movie)

		err := svc.DeleteMovie(ctx, movie.ID)
		if err != nil {
			t.Fatalf("DeleteMovie() error = %v", err)
		}

		_, err = repo.GetByID(ctx, movie.ID)
		if !errors.Is(err, ErrMovieNotFound) {
			t.Errorf("DeleteMovie() movie still exists")
		}
	})
}

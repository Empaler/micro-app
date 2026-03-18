package movies

import (
	"context"
	"errors"

	"movie-api/internal/domain"
)

var ErrMovieNotFound = errors.New("movie not found")

type MovieService struct {
	repo shared.Repository[Movie]
}

func NewMovieService(repo shared.Repository[Movie]) *MovieService {
	return &MovieService{repo: repo}
}

func (s *MovieService) CreateMovie(ctx context.Context, movie *Movie) error {
	if err := movie.Validate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, movie)
}

func (s *MovieService) GetMovie(ctx context.Context, id int64) (*Movie, error) {
	movie, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrMovieNotFound) {
			return nil, ErrMovieNotFound
		}
		return nil, err
	}
	return movie, nil
}

func (s *MovieService) GetAllMovies(ctx context.Context) ([]Movie, error) {
	return s.repo.GetAll(ctx)
}

func (s *MovieService) UpdateMovie(ctx context.Context, movie *Movie) error {
	if err := movie.Validate(); err != nil {
		return err
	}

	existing, err := s.repo.GetByID(ctx, movie.ID)
	if err != nil {
		if errors.Is(err, ErrMovieNotFound) {
			return ErrMovieNotFound
		}
		return err
	}

	existing.Title = movie.Title
	existing.Year = movie.Year
	existing.Type = movie.Type
	existing.Resolution = movie.Resolution
	existing.Actors = movie.Actors
	existing.Rating = movie.Rating
	existing.IsAdult = movie.IsAdult

	return s.repo.Update(ctx, existing)
}

func (s *MovieService) DeleteMovie(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, ErrMovieNotFound) {
			return ErrMovieNotFound
		}
		return err
	}
	return nil
}

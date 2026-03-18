package books

import (
	"context"
	"errors"

	"movie-api/internal/domain"
)

var ErrBookNotFound = errors.New("book not found")

type BookService struct {
	repo shared.Repository[Book]
}

func NewBookService(repo shared.Repository[Book]) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) CreateBook(ctx context.Context, book *Book) error {
	if err := book.Validate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, book)
}

func (s *BookService) GetBook(ctx context.Context, id int64) (*Book, error) {
	book, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			return nil, ErrBookNotFound
		}
		return nil, err
	}
	return book, nil
}

func (s *BookService) GetAllBooks(ctx context.Context) ([]Book, error) {
	return s.repo.GetAll(ctx)
}

func (s *BookService) UpdateBook(ctx context.Context, book *Book) error {
	if err := book.Validate(); err != nil {
		return err
	}

	existing, err := s.repo.GetByID(ctx, book.ID)
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			return ErrBookNotFound
		}
		return err
	}

	existing.Title = book.Title
	existing.Author = book.Author
	existing.ReleaseYear = book.ReleaseYear
	existing.Rating = book.Rating

	return s.repo.Update(ctx, existing)
}

func (s *BookService) DeleteBook(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			return ErrBookNotFound
		}
		return err
	}
	return nil
}

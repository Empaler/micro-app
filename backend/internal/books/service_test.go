package books

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockBookRepo struct {
	books  map[int64]*Book
	nextID int64
}

func newMockBookRepo() *mockBookRepo {
	return &mockBookRepo{books: make(map[int64]*Book), nextID: 1}
}

func (m *mockBookRepo) Create(ctx context.Context, book *Book) error {
	if book.ID == 0 {
		book.ID = m.nextID
		m.nextID++
	}
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()
	m.books[book.ID] = book
	return nil
}

func (m *mockBookRepo) GetByID(ctx context.Context, id int64) (*Book, error) {
	if book, ok := m.books[id]; ok {
		return book, nil
	}
	return nil, ErrBookNotFound
}

func (m *mockBookRepo) GetAll(ctx context.Context) ([]Book, error) {
	result := make([]Book, 0, len(m.books))
	for _, book := range m.books {
		result = append(result, *book)
	}
	return result, nil
}

func (m *mockBookRepo) Update(ctx context.Context, book *Book) error {
	if _, ok := m.books[book.ID]; !ok {
		return ErrBookNotFound
	}
	book.UpdatedAt = time.Now()
	m.books[book.ID] = book
	return nil
}

func (m *mockBookRepo) Delete(ctx context.Context, id int64) error {
	if _, ok := m.books[id]; !ok {
		return ErrBookNotFound
	}
	delete(m.books, id)
	return nil
}

func TestBookService_CreateBook(t *testing.T) {
	tests := []struct {
		name    string
		book    Book
		wantErr bool
	}{
		{
			name: "valid book",
			book: Book{
				Title:       "The Go Programming Language",
				Author:      "Donovan",
				ReleaseYear: 2015,
				Rating:      9.5,
			},
			wantErr: false,
		},
		{
			name: "missing title",
			book: Book{
				Title:       "",
				Author:      "Author",
				ReleaseYear: 2020,
				Rating:      7.0,
			},
			wantErr: true,
		},
		{
			name: "missing author",
			book: Book{
				Title:       "Some Book",
				Author:      "",
				ReleaseYear: 2020,
				Rating:      7.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockBookRepo()
			svc := NewBookService(repo)

			err := svc.CreateBook(context.Background(), &tt.book)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBookService_GetBook(t *testing.T) {
	ctx := context.Background()
	repo := newMockBookRepo()
	svc := NewBookService(repo)

	t.Run("book not found", func(t *testing.T) {
		_, err := svc.GetBook(ctx, 999)
		if !errors.Is(err, ErrBookNotFound) {
			t.Errorf("GetBook() error = %v, want %v", err, ErrBookNotFound)
		}
	})

	t.Run("get existing book", func(t *testing.T) {
		book := &Book{
			Title:       "Test Book",
			Author:      "Test Author",
			ReleaseYear: 2020,
			Rating:      7.5,
		}
		repo.Create(ctx, book)

		got, err := svc.GetBook(ctx, book.ID)
		if err != nil {
			t.Fatalf("GetBook() error = %v", err)
		}
		if got.Title != book.Title {
			t.Errorf("GetBook() title = %v, want %v", got.Title, book.Title)
		}
	})
}

func TestBookService_GetAllBooks(t *testing.T) {
	ctx := context.Background()
	repo := newMockBookRepo()
	svc := NewBookService(repo)

	repo.Create(ctx, &Book{
		Title:       "Book 1",
		Author:      "Author 1",
		ReleaseYear: 2020,
		Rating:      7.0,
	})
	repo.Create(ctx, &Book{
		Title:       "Book 2",
		Author:      "Author 2",
		ReleaseYear: 2021,
		Rating:      8.0,
	})

	books, err := svc.GetAllBooks(ctx)
	if err != nil {
		t.Fatalf("GetAllBooks() error = %v", err)
	}
	if len(books) != 2 {
		t.Errorf("GetAllBooks() got %d books, want 2", len(books))
	}
}

func TestBookService_UpdateBook(t *testing.T) {
	ctx := context.Background()
	repo := newMockBookRepo()
	svc := NewBookService(repo)

	t.Run("update non-existent book", func(t *testing.T) {
		err := svc.UpdateBook(ctx, &Book{
			ID:          999,
			Title:       "Updated",
			Author:      "Author",
			ReleaseYear: 2020,
			Rating:      7.0,
		})
		if !errors.Is(err, ErrBookNotFound) {
			t.Errorf("UpdateBook() error = %v, want %v", err, ErrBookNotFound)
		}
	})

	t.Run("update existing book", func(t *testing.T) {
		book := &Book{
			Title:       "Original",
			Author:      "Author",
			ReleaseYear: 2020,
			Rating:      7.0,
		}
		repo.Create(ctx, book)

		book.Title = "Updated Title"
		err := svc.UpdateBook(ctx, book)
		if err != nil {
			t.Fatalf("UpdateBook() error = %v", err)
		}

		got, _ := repo.GetByID(ctx, book.ID)
		if got.Title != "Updated Title" {
			t.Errorf("UpdateBook() title = %v, want Updated Title", got.Title)
		}
	})

	t.Run("update with invalid data", func(t *testing.T) {
		book := &Book{
			Title:       "Test",
			Author:      "Author",
			ReleaseYear: 2020,
			Rating:      7.0,
		}
		repo.Create(ctx, book)

		book.Title = ""
		err := svc.UpdateBook(ctx, book)
		if err != ErrBookTitleRequired {
			t.Errorf("UpdateBook() error = %v, want %v", err, ErrBookTitleRequired)
		}
	})
}

func TestBookService_DeleteBook(t *testing.T) {
	ctx := context.Background()
	repo := newMockBookRepo()
	svc := NewBookService(repo)

	t.Run("delete non-existent book", func(t *testing.T) {
		err := svc.DeleteBook(ctx, 999)
		if !errors.Is(err, ErrBookNotFound) {
			t.Errorf("DeleteBook() error = %v, want %v", err, ErrBookNotFound)
		}
	})

	t.Run("delete existing book", func(t *testing.T) {
		book := &Book{
			Title:       "Test",
			Author:      "Author",
			ReleaseYear: 2020,
			Rating:      7.0,
		}
		repo.Create(ctx, book)

		err := svc.DeleteBook(ctx, book.ID)
		if err != nil {
			t.Fatalf("DeleteBook() error = %v", err)
		}

		_, err = repo.GetByID(ctx, book.ID)
		if !errors.Is(err, ErrBookNotFound) {
			t.Errorf("DeleteBook() book still exists")
		}
	})
}

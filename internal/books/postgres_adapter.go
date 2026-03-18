package books

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type PostgresAdapter struct {
	db *sqlx.DB
}

func NewPostgresAdapter(db *sqlx.DB) *PostgresAdapter {
	return &PostgresAdapter{db: db}
}

func (a *PostgresAdapter) Create(ctx context.Context, book *Book) error {
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()
	query := `
		INSERT INTO books (title, author, release_year, rating, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	return a.db.QueryRowContext(ctx, query,
		book.Title, book.Author, book.ReleaseYear, book.Rating, book.CreatedAt, book.UpdatedAt,
	).Scan(&book.ID)
}

func (a *PostgresAdapter) GetByID(ctx context.Context, id int64) (*Book, error) {
	var book Book
	query := `SELECT id, title, author, release_year, rating, created_at, updated_at FROM books WHERE id = $1`
	err := a.db.GetContext(ctx, &book, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrBookNotFound
		}
		return nil, err
	}
	return &book, nil
}

func (a *PostgresAdapter) GetAll(ctx context.Context) ([]Book, error) {
	var books []Book
	query := `SELECT id, title, author, release_year, rating, created_at, updated_at FROM books ORDER BY created_at DESC`
	err := a.db.SelectContext(ctx, &books, query)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (a *PostgresAdapter) Update(ctx context.Context, book *Book) error {
	book.UpdatedAt = time.Now()
	query := `
		UPDATE books 
		SET title = $1, author = $2, release_year = $3, rating = $4, updated_at = $5
		WHERE id = $6`

	result, err := a.db.ExecContext(ctx, query,
		book.Title, book.Author, book.ReleaseYear, book.Rating, book.UpdatedAt, book.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrBookNotFound
	}

	return nil
}

func (a *PostgresAdapter) Delete(ctx context.Context, id int64) error {
	result, err := a.db.ExecContext(ctx, "DELETE FROM books WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrBookNotFound
	}
	return nil
}

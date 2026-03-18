package movies

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

func (a *PostgresAdapter) Create(ctx context.Context, movie *Movie) error {
	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()
	query := `
		INSERT INTO movies (title, year, type, resolution, actors, rating, is_adult, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	return a.db.QueryRowContext(ctx, query,
		movie.Title, movie.Year, movie.Type, movie.Resolution,
		movie.Actors, movie.Rating, movie.IsAdult, movie.CreatedAt, movie.UpdatedAt,
	).Scan(&movie.ID)
}

func (a *PostgresAdapter) GetByID(ctx context.Context, id int64) (*Movie, error) {
	var movie Movie
	query := `SELECT id, title, year, type, resolution, actors, rating, is_adult, created_at, updated_at FROM movies WHERE id = $1`
	err := a.db.GetContext(ctx, &movie, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrMovieNotFound
		}
		return nil, err
	}
	return &movie, nil
}

func (a *PostgresAdapter) GetAll(ctx context.Context) ([]Movie, error) {
	var movies []Movie
	query := `SELECT id, title, year, type, resolution, actors, rating, is_adult, created_at, updated_at FROM movies ORDER BY created_at DESC`
	err := a.db.SelectContext(ctx, &movies, query)
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (a *PostgresAdapter) Update(ctx context.Context, movie *Movie) error {
	movie.UpdatedAt = time.Now()
	query := `
		UPDATE movies 
		SET title = $1, year = $2, type = $3, resolution = $4, actors = $5, rating = $6, is_adult = $7, updated_at = $8
		WHERE id = $9`

	result, err := a.db.ExecContext(ctx, query,
		movie.Title, movie.Year, movie.Type, movie.Resolution,
		movie.Actors, movie.Rating, movie.IsAdult, movie.UpdatedAt, movie.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrMovieNotFound
	}

	return nil
}

func (a *PostgresAdapter) Delete(ctx context.Context, id int64) error {
	result, err := a.db.ExecContext(ctx, "DELETE FROM movies WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrMovieNotFound
	}
	return nil
}

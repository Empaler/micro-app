package books

import (
	"errors"
	"time"
)

var (
	ErrBookTitleRequired  = errors.New("title is required")
	ErrBookAuthorRequired = errors.New("author is required")
	ErrBookInvalidYear    = errors.New("year must be between 1000 and current year")
	ErrBookInvalidRating  = errors.New("rating must be between 0 and 10")
)

type Book struct {
	ID          int64     `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Author      string    `json:"author" db:"author"`
	ReleaseYear int       `json:"releaseYear" db:"release_year"`
	Rating      float64   `json:"rating" db:"rating"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

func (b *Book) Validate() error {
	if b.Title == "" {
		return ErrBookTitleRequired
	}
	if b.Author == "" {
		return ErrBookAuthorRequired
	}

	currentYear := time.Now().Year()
	if b.ReleaseYear < 1000 || b.ReleaseYear > currentYear {
		return ErrBookInvalidYear
	}

	if b.Rating < 0 || b.Rating > 10 {
		return ErrBookInvalidRating
	}

	return nil
}

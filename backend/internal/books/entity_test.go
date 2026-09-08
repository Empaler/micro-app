package books

import (
	"testing"
	"time"
)

func TestBook_Validate(t *testing.T) {
	currentYear := time.Now().Year()

	tests := []struct {
		name    string
		book    Book
		wantErr error
	}{
		{
			name: "valid book",
			book: Book{
				Title:       "The Go Programming Language",
				Author:      "Alan A. A. Donovan",
				ReleaseYear: 2015,
				Rating:      9.5,
			},
			wantErr: nil,
		},
		{
			name: "missing title",
			book: Book{
				Title:       "",
				Author:      "Author Name",
				ReleaseYear: 2020,
				Rating:      7.0,
			},
			wantErr: ErrBookTitleRequired,
		},
		{
			name: "missing author",
			book: Book{
				Title:       "Some Book",
				Author:      "",
				ReleaseYear: 2020,
				Rating:      7.0,
			},
			wantErr: ErrBookAuthorRequired,
		},
		{
			name: "year too early",
			book: Book{
				Title:       "Ancient Book",
				Author:      "Ancient Author",
				ReleaseYear: 999,
				Rating:      5.0,
			},
			wantErr: ErrBookInvalidYear,
		},
		{
			name: "year too late",
			book: Book{
				Title:       "Future Book",
				Author:      "Future Author",
				ReleaseYear: currentYear + 1,
				Rating:      5.0,
			},
			wantErr: ErrBookInvalidYear,
		},
		{
			name: "rating too low",
			book: Book{
				Title:       "Test Book",
				Author:      "Test Author",
				ReleaseYear: 2020,
				Rating:      -0.1,
			},
			wantErr: ErrBookInvalidRating,
		},
		{
			name: "rating too high",
			book: Book{
				Title:       "Test Book",
				Author:      "Test Author",
				ReleaseYear: 2020,
				Rating:      10.1,
			},
			wantErr: ErrBookInvalidRating,
		},
		{
			name: "valid with boundary year",
			book: Book{
				Title:       "Old Book",
				Author:      "Author",
				ReleaseYear: 1000,
				Rating:      5.0,
			},
			wantErr: nil,
		},
		{
			name: "valid rating boundary 0",
			book: Book{
				Title:       "Bad Book",
				Author:      "Author",
				ReleaseYear: 2020,
				Rating:      0,
			},
			wantErr: nil,
		},
		{
			name: "valid rating boundary 10",
			book: Book{
				Title:       "Great Book",
				Author:      "Author",
				ReleaseYear: 2020,
				Rating:      10,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.book.Validate()
			if err != tt.wantErr {
				t.Errorf("Book.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

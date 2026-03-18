package movies

import (
	"testing"
	"time"
)

func TestMovie_Validate(t *testing.T) {
	currentYear := time.Now().Year()

	tests := []struct {
		name    string
		movie   Movie
		wantErr error
	}{
		{
			name: "valid movie",
			movie: Movie{
				Title:      "The Matrix",
				Year:       1999,
				Type:       MovieTypeMovie,
				Resolution: ResolutionFHD,
				Actors:     "Keanu Reeves",
				Rating:     8.7,
				IsAdult:    false,
			},
			wantErr: nil,
		},
		{
			name: "valid series",
			movie: Movie{
				Title:      "Breaking Bad",
				Year:       2008,
				Type:       MovieTypeSeries,
				Resolution: ResolutionHD,
				Rating:     9.5,
			},
			wantErr: nil,
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
			wantErr: ErrTitleRequired,
		},
		{
			name: "year too early",
			movie: Movie{
				Title:      "Test",
				Year:       1887,
				Type:       MovieTypeMovie,
				Resolution: ResolutionHD,
				Rating:     5.0,
			},
			wantErr: ErrInvalidYear,
		},
		{
			name: "year too late",
			movie: Movie{
				Title:      "Test",
				Year:       currentYear + 1,
				Type:       MovieTypeMovie,
				Resolution: ResolutionHD,
				Rating:     5.0,
			},
			wantErr: ErrInvalidYear,
		},
		{
			name: "rating too low",
			movie: Movie{
				Title:      "Test",
				Year:       2000,
				Type:       MovieTypeMovie,
				Resolution: ResolutionHD,
				Rating:     -0.1,
			},
			wantErr: ErrInvalidRating,
		},
		{
			name: "rating too high",
			movie: Movie{
				Title:      "Test",
				Year:       2000,
				Type:       MovieTypeMovie,
				Resolution: ResolutionHD,
				Rating:     10.1,
			},
			wantErr: ErrInvalidRating,
		},
		{
			name: "invalid type",
			movie: Movie{
				Title:      "Test",
				Year:       2000,
				Type:       "documentary",
				Resolution: ResolutionHD,
				Rating:     5.0,
			},
			wantErr: ErrInvalidType,
		},
		{
			name: "invalid resolution",
			movie: Movie{
				Title:      "Test",
				Year:       2000,
				Type:       MovieTypeMovie,
				Resolution: "8K",
				Rating:     5.0,
			},
			wantErr: ErrInvalidResolution,
		},
		{
			name: "adult content not allowed",
			movie: Movie{
				Title:      "Test",
				Year:       2000,
				Type:       MovieTypeMovie,
				Resolution: ResolutionHD,
				Rating:     5.0,
				IsAdult:    true,
			},
			wantErr: ErrAdultContentNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.movie.Validate()
			if err != tt.wantErr {
				t.Errorf("Movie.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

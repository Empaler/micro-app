package movies

import (
	"errors"
	"time"
)

var (
	ErrAdultContentNotAllowed = errors.New("adult content is not allowed")
	ErrInvalidYear            = errors.New("year must be between 1888 and current year")
	ErrInvalidRating          = errors.New("rating must be between 0 and 10")
	ErrInvalidType            = errors.New("type must be 'movie' or 'series'")
	ErrInvalidResolution      = errors.New("resolution must be 'SD', 'HD', 'FHD', or '4K'")
	ErrTitleRequired          = errors.New("title is required")
)

type MovieType string

const (
	MovieTypeMovie  MovieType = "movie"
	MovieTypeSeries MovieType = "series"
)

type Resolution string

const (
	ResolutionSD  Resolution = "SD"
	ResolutionHD  Resolution = "HD"
	ResolutionFHD Resolution = "FHD"
	Resolution4K  Resolution = "4K"
)

type Movie struct {
	ID         int64      `json:"id" db:"id"`
	Title      string     `json:"title" db:"title"`
	Year       int        `json:"year" db:"year"`
	Type       MovieType  `json:"type" db:"type"`
	Resolution Resolution `json:"resolution" db:"resolution"`
	Actors     string     `json:"actors" db:"actors"`
	Rating     float64    `json:"rating" db:"rating"`
	IsAdult    bool       `json:"isAdult" db:"is_adult"`
	CreatedAt  time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time  `json:"updatedAt" db:"updated_at"`
}

func (m *Movie) Validate() error {
	if m.Title == "" {
		return ErrTitleRequired
	}

	currentYear := time.Now().Year()
	if m.Year < 1888 || m.Year > currentYear {
		return ErrInvalidYear
	}

	if m.Rating < 0 || m.Rating > 10 {
		return ErrInvalidRating
	}

	if m.Type != MovieTypeMovie && m.Type != MovieTypeSeries {
		return ErrInvalidType
	}

	if m.Resolution != ResolutionSD && m.Resolution != ResolutionHD &&
		m.Resolution != ResolutionFHD && m.Resolution != Resolution4K {
		return ErrInvalidResolution
	}

	if m.IsAdult {
		return ErrAdultContentNotAllowed
	}

	return nil
}

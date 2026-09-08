-- Create movies table
CREATE TABLE IF NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    year INTEGER NOT NULL,
    type VARCHAR(20) NOT NULL,
    resolution VARCHAR(10) NOT NULL,
    actors TEXT,
    rating REAL NOT NULL,
    is_adult BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_year CHECK (year >= 1888 AND year <= EXTRACT(YEAR FROM CURRENT_DATE)),
    CONSTRAINT chk_type CHECK (type IN ('movie', 'series')),
    CONSTRAINT chk_resolution CHECK (resolution IN ('SD', 'HD', 'FHD', '4K')),
    CONSTRAINT chk_rating CHECK (rating >= 0 AND rating <= 10)
);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_movies_created_at ON movies(created_at DESC);

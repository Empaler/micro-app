-- Create books table
CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    release_year INTEGER NOT NULL,
    rating REAL NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_release_year CHECK (release_year >= 1000 AND release_year <= EXTRACT(YEAR FROM CURRENT_DATE)),
    CONSTRAINT chk_rating CHECK (rating >= 0 AND rating <= 10)
);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_books_created_at ON books(created_at DESC);

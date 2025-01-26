CREATE TABLE books (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    isbn TEXT UNIQUE NOT NULL,
    author TEXT NOT NULL,
    number_of_pages INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
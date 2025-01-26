CREATE TABLE book_items (
    id UUID PRIMARY KEY,
    book_id UUID NOT NULL,
    location TEXT NOT NULL,
    condition TEXT,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_book FOREIGN KEY(book_id) REFERENCES books(id) ON DELETE CASCADE
);
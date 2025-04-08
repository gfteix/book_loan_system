CREATE TABLE loans (
    id UUID PRIMARY KEY,
    book_item_id UUID NOT NULL,
    user_id UUID NOT NULL,
    status TEXT NOT NULL,
    loan_date TIMESTAMP,
    expiring_date TIMESTAMP,
    return_date TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_book_item_id FOREIGN KEY(book_item_id) REFERENCES book_copies(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
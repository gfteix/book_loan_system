INSERT INTO book_copies (id, book_id, location, condition, status) 
SELECT gen_random_uuid(), id, 'Shelf A', 'New', 'Available' 
FROM books;
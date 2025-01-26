INSERT INTO loans (id, book_item_id, user_id, status, loan_date, expiring_date) 
VALUES 
    (gen_random_uuid(), 
        (SELECT id FROM book_items LIMIT 1 OFFSET 0), 
        (SELECT id FROM users LIMIT 1), 
        'Active', 
        CURRENT_TIMESTAMP, 
        CURRENT_DATE + INTERVAL '1 day'),
    (gen_random_uuid(), 
        (SELECT id FROM book_items LIMIT 1 OFFSET 1), 
        (SELECT id FROM users LIMIT 1), 
        'Active', 
        CURRENT_TIMESTAMP, 
        CURRENT_DATE + INTERVAL '1 day');
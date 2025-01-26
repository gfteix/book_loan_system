INSERT INTO books (id, title, description, isbn, author, number_of_pages) 
VALUES 
    (gen_random_uuid(), 'Crime and Punishment', 'A psychological novel by Dostoyevsky', '9780143058144', 'Fyodor Dostoyevsky', 671),
    (gen_random_uuid(), 'The Brothers Karamazov', 'A novel by Dostoyevsky about faith and doubt', '9780140449242', 'Fyodor Dostoyevsky', 824),
    (gen_random_uuid(), 'War and Peace', 'A historical novel by Tolstoy', '9780199232765', 'Leo Tolstoy', 1392),
    (gen_random_uuid(), 'Anna Karenina', 'A romantic tragedy by Tolstoy', '9780143035008', 'Leo Tolstoy', 864),
    (gen_random_uuid(), 'The Metamorphosis', 'A novella by Kafka', '9780805210576', 'Franz Kafka', 104);
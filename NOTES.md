
## Database

- To create a postgreSQL container

```

docker run --name postgres_docker -p 5432:5432 -e POSTGRES_PASSWORD=mypassword -e POSTGRES_USER=postgres -e POSTGRES_DB=library  -d postgres

```


- To stop and remove a container (with the volume)

```
docker container stop postgres_docker && docker container remove postgres_docker && docker volume prune 
```


- To go inside the container:

`docker exec -it postgres_docker psql -U postgres -d library` (where “root” is the username for MySQL database.)

- To list all tables in the current schema:


```
\d
```

SQL

```
CREATE TABLE users (
    id UUID PRIMARY KEY,     
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE books (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    isbn TEXT UNIQUE NOT NULL,
    author TEXT NOT NULL,
    number_of_pages INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE book_items (
    id UUID PRIMARY KEY,
    book_id UUID NOT NULL,
    location TEXT NOT NULL,
    condition TEXT,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_book FOREIGN KEY(book_id) REFERENCES books(id) ON DELETE CASCADE
);

CREATE TABLE loans (
    id UUID PRIMARY KEY,
    book_item_id UUID NOT NULL,
    user_id UUID NOT NULL,
    status TEXT NOT NULL,
    loan_date TIMESTAMP,
    expiring_date TIMESTAMP,
    return_date TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_book_item_id FOREIGN KEY(book_item_id) REFERENCES book_items(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

```


---

## Requests

```
curl -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{
  "name": "John",
  "email": "john@example.com"
}' -v

```

```
curl http://localhost:8080/users/528a1dbc-d391-46e3-b818-6cf78e4344d2
```

```

curl -X POST http://localhost:8080/books \
-H "Content-Type: application/json" \
-d '{
  "title": "title",
  "description": "some description",
  "isbn": "1",
  "author": "author",
  "numberOfPages": 100
}' -v

```

```
curl "http://localhost:8080/books?title=sometitle"
```

```
curl http://localhost:8080/books/cca29657-a87d-4300-a4b4-a3163a054872
```


```
curl -X POST http://localhost:8080/books/{id}/items \
-H "Content-Type: application/json" \
-d '{
  "bookId": "cca29657-a87d-4300-a4b4-a3163a054872",
  "status": "available",
  "condition":"good",
  "location": "section b"
}' -v

```

```
curl http://localhost:8080/books/cca29657-a87d-4300-a4b4-a3163a054872/items
```

- Create Loans

```
curl -X POST http://localhost:8080/loans \
-H "Content-Type: application/json" \
-d '{
  "userId": "528a1dbc-d391-46e3-b818-6cf78e4344d2",
  "status": "active",
  "bookItemId":"15caa834-4c8b-4c75-833f-bbad805e8a3c",
  "loanDate": "2025-01-12T15:30:00Z",
  "expiringDate": "2025-01-14T15:30:00Z"
}' -v

```

- Get Loans

curl "http://localhost:8080/loans?userId=528a1dbc-d391-46e3-b818-6cf78e4344d2"

---

Email handler should receive an event type along with the userId/loanId

Notifications:
- should send reminder to return the book; reminder should be send 2 days before the expected return date
- should send alert on the day of the expiring date if not returned yet

Payloads structure:

```
{
  "source": string, // file source or api source endpoint
  "time": string, // time of event generation
  "event-id": string, // uuid
  "type": string // "loan-expiring" || "loan-ended"
  "payload": {
    "userId": string,
    "loanId": string,
  }
}
```

- Use DLQ if there are failures when sending email

---

## RabbitMQ

Docker command:

```
# latest RabbitMQ 4.0.x
docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:4.0-management
```

- Port 5672: Enables application-level communication with RabbitMQ (the main messaging functionality).

- Port 15672: Allows you to monitor and manage RabbitMQ using a web browser.
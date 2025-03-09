
# Book Loan System

Backend System to manage book loans.

## Functionalities

- Create user
- Get user
- Get users
- Get List of Books (with filters)
- Get Book details
- Get Book item details
- Lend a book (item)
- Return a book (item)
- Send email to client when loan is about to expire

Database Diagram

![database modeling](book_loan_system-DB.drawio.png "Title")

System Design

![system design](book_loan_system-system_design.drawio.png "Title")

## Technologies

- Go
- Docker
- RabbitMQ

## TODO

- [X] Build API
- [X] Build loan expiring job
- [X] Build email handler
- [X] Add database migrations
- [X] Dockerfiles
- [X] Docker Compose
- [ ] Swagger Docs
- [ ] Integrate with Prometheus for metrics?
- [ ] Rename book_items to book_copies (db, code, db diagram and NOTES.md)
---

## How to run

- Create a `.env` filed based on the `.env.example`

- build the infraestructure and execute the project by running

`docker compose --env-file .env build --no-cache && docker compose --env-file .env up -d --force-recreate`




## API


- To create an user
```
curl -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{
  "name": "John",
  "email": "john@example.com"
}' -v

```

- To get an user
```
curl http://localhost:8080/users/{id}
```

- To create a book

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

- To search a book by title

```
curl "http://localhost:8080/books?title=sometitle"
```

- To get a book by ID

```
curl http://localhost:8080/books/cca29657-a87d-4300-a4b4-a3163a054872
```

- To create a book item

```
curl -X POST http://localhost:8080/books/{id}/items \
-H "Content-Type: application/json" \
-d '{
  "bookId": "3471807e-1c3b-4b27-b397-8f9123e6a6f0",
  "status": "available",
  "condition":"good",
  "location": "section b"
}' -v

```

- To get book items of a book

```
curl http://localhost:8080/books/3471807e-1c3b-4b27-b397-8f9123e6a6f0/items
```

- To create a loan

```
curl -X POST http://localhost:8080/loans \
-H "Content-Type: application/json" \
-d '{
  "userId": "2b0e169b-55d9-4356-ba44-3aa23dd9b2a0",
  "status": "active",
  "bookItemId":"36fbab72-3a61-46f0-a211-7619bc2916c5",
  "loanDate": "2025-01-26T15:30:00Z",
  "expiringDate": "2025-01-27T15:30:00Z"
}' -v

```

- To Get Loans

```
curl "http://localhost:8080/loans?userId={id}"
```

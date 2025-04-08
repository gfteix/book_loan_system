# Book Loan System

A backend system to manage book loans.

## Features

- Create, retrieve, and list users
- Retrieve books with filters
- Retrieve book details and book item details
- Lend and return book items
- Notify users via email when a loan is about to expire

## System Overview

### Database Diagram
![Database Modeling](book_loan_system-DB.drawio.png "Database Schema")

### System Design
![System Design](book_loan_system-system_design.drawio.png "System Architecture")

## Technologies Used

- **Go** - Backend development
- **Docker** - Containerization
- **RabbitMQ** - Message queuing

## TODO

- [X] Build API
- [X] Build loan expiring job
- [X] Build email handler
- [X] Add database migrations
- [X] Dockerfiles
- [X] Docker Compose
- [X] Swagger Docs
- [X] Rename book_copies to book_copies (db, code, db diagram and NOTES.md)
- [ ] Integrate with Prometheus for metrics?
---


## Getting Started

### Setup and Run

1. Create a `.env` file based on `.env.example`.
2. Build and start the infrastructure:

   ```sh
   docker compose --env-file .env build --no-cache && docker compose --env-file .env up -d --force-recreate
   ```


### API Documentation

Once the API is running, you can access the Swagger docs at:

[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

Alternatively, you can open the HTML documentation manually from the `docs/` folder in a browser.

## API Endpoints

### User Management

#### Create a User
```sh
curl -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{
  "name": "John",
  "email": "john@example.com"
}' -v
```

#### Get a User by ID
```sh
curl http://localhost:8080/users/{id}
```

### Book Management

#### Create a Book
```sh
curl -X POST http://localhost:8080/books \
-H "Content-Type: application/json" \
-d '{
  "title": "Book Title",
  "description": "A detailed description",
  "isbn": "123456789",
  "author": "Author Name",
  "numberOfPages": 250
}' -v
```

#### Search Books by Title
```sh
curl "http://localhost:8080/books?title=example"
```

#### Get a Book by ID
```sh
curl http://localhost:8080/books/{book_id}
```

### Book Item Management

#### Create a Book Item
```sh
curl -X POST http://localhost:8080/books/{book_id}/items \
-H "Content-Type: application/json" \
-d '{
  "bookId": "book_uuid",
  "status": "available",
  "condition": "good",
  "location": "Section B"
}' -v
```

#### Get Book Items for a Book
```sh
curl http://localhost:8080/books/{book_id}/items
```

### Loan Management

#### Create a Loan
```sh
curl -X POST http://localhost:8080/loans \
-H "Content-Type: application/json" \
-d '{
  "userId": "user_uuid",
  "status": "active",
  "bookCopyId": "book_item_uuid",
  "loanDate": "2025-01-26T15:30:00Z",
  "expiringDate": "2025-01-27T15:30:00Z"
}' -v
```

#### Get Loans by User ID
```sh
curl "http://localhost:8080/loans?userId={user_id}"
```


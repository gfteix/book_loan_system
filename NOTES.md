
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
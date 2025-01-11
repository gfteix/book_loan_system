
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
```


---

## Requests

```
curl -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{
  "name": "John",
  "email": "john@example.com",
}' -v

```
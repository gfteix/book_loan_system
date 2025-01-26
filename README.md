
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
- [ ] Add Rate Limit to API
- [X] Add database migrations
- [ ] Dockerfiles
- [ ] Docker Compose (in progress)
- [ ] Swagger Docs
- [ ] Integrate with Prometheus for metrics


---

## How to run

- Make sure the infraestructure is up

`docker compose up`


### API

`make api-run`

### Email Handler

The handler will be listening to messages on the rabbitmq queue, when a new message arrives it sends an email.
It is possible to see the emails sent on http://localhost:8025/

`go run cmd/emails/main.go`

### Reminder Job 

It checks for existing loans that expires today or will expire in the next day and sends a message to the queue

`go run cmd/reminders/main.go`


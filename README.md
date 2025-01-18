
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
- [ ] Build loan expiring job (developing)
- [ ] Build email handler (developing)
- [ ] Add Rate Limit to API
- [ ] Add database migrations
- [ ] Dockerfiles and Docker Compose
- [ ] Swagger Docs
- [ ] Integrate with Prometheus for metrics


---

## How to run

### API

`make api-run`

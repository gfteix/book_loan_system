api-build:
	@go build -o bin/api cmd/api/main.go

api-run: api-build
	@./bin/api

test:
	@go test -v ./... -cover

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

reminder-build:
	@go build -o bin/reminder cmd/loan-reminder/main.go

reminder-run: reminder-build
	@./bin/reminder
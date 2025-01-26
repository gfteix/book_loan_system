api-build:
	@go build -o bin/api cmd/api/main.go

api-run: api-build
	@./bin/api

reminders-build:
	@go build -o bin/reminders cmd/reminders/main.go

reminders-run: reminders-build
	@./bin/reminders

emails-build:
	@go build -o bin/emails cmd/emails/main.go

emails-run: emails-build
	@./bin/emails

test:
	@go test -v ./... -cover

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

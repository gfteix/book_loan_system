FROM golang:1.22.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/reminders ./cmd/reminders

FROM gcr.io/distroless/base-debian10

COPY --from=builder /bin/reminders /bin/reminders

EXPOSE 8181

CMD ["/bin/reminders"]

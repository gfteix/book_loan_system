FROM golang:1.23.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify || (echo "Go mod failed" && cat go.sum && exit 1)

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/emails ./cmd/emails

FROM gcr.io/distroless/base-debian10

COPY --from=builder /bin/emails /bin/emails

EXPOSE 6060

CMD ["/bin/emails"]

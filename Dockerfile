# build env/setup
FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# install dependencies for docs and db migrations
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN make docs
RUN make goose-up
RUN make build

# runtime
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]


# build env/setup
FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make docs
RUN make build

# runtime
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]


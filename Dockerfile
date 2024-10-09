# Stage 1: Build the Go binary
FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o /app/binary

# Stage 2: Create a minimal runtime image
FROM alpine:3.14

RUN mkdir -p /app/uploads /app/logs

COPY --from=builder /app/binary /app/binary

ENTRYPOINT ["/app/binary"]

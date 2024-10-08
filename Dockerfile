# Stage 1: Build the Go binary
FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o /app/binary

# Stage 2: Create a minimal runtime image
FROM alpine:3.14

COPY --from=builder /app/binary /app/binary

# Buat direktori untuk uploads dan logs
RUN mkdir -p /app/uploads /app/logs

# Mendefinisikan volume untuk uploads dan logs
VOLUME ["/app/uploads"]
VOLUME ["/app/logs"]

ENTRYPOINT ["/app/binary"]

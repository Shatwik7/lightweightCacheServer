# Build Stage
FROM golang:1.23.5-alpine AS builder

WORKDIR /app

# Install binutils for 'strip' command
RUN apk add --no-cache binutils

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main . && \
    strip main  # Reduce binary size

# Final Minimal Image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 2345

CMD ["./main"]

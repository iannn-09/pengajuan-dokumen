# Stage 1: Build binary
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy dependency definition
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api/main.go

# Stage 2: Final minimal runtime image
FROM alpine:latest  

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .env

# Expose API port
EXPOSE 8080

CMD ["./main"]

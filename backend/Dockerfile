# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git for any go module fetch
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Cache buster - forces rebuild when commit SHA changes
ARG CACHEBUST=1

# Copy source code
COPY . .

# Build the Go API application binary statically
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api/main.go

# Production stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates & tzdata for HTTPS & Jakarta timezone
RUN apk --no-cache add ca-certificates tzdata

# Set timezone
ENV TZ=Asia/Jakarta

# Copy binary from builder
COPY --from=builder /app/main .

# Create uploads directory for document files & avatars
RUN mkdir -p /app/uploads

# Expose API port 8080
EXPOSE 8080

CMD ["./main"]

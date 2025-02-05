# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

# 添加构建依赖
RUN apk add --no-cache gcc musl-dev

# Download dependencies
RUN go mod download

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# 添加 CA 证书，用于 HTTPS 请求
RUN apk add --no-cache ca-certificates

# 添加非 root 用户
RUN adduser -D appuser
USER appuser

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# 添加环境变量
ENV PORT=80
ENV GIN_MODE=release

EXPOSE 80

CMD ["./main"]
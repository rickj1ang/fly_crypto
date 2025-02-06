# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM ubuntu:20.04

WORKDIR /app

# 安装 CA 证书
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# 创建非 root 用户
RUN useradd -r -u 1001 -g root appuser
USER appuser

COPY --from=builder /app/main .

ENV PORT=80
ENV GIN_MODE=release

EXPOSE 80

CMD ["./main"]

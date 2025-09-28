# 后端 Go 服务 Dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main ./main
COPY roles.json ./roles.json
COPY config ./config
COPY models ./models
COPY services ./services
COPY handlers ./handlers
EXPOSE 8080
ENV PORT=8080
CMD ["./main"]

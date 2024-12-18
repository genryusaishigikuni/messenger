FROM golang:1.20-alpine as builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o presence-service ./cmd/presence

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/presence-service .
ENV AUTH_SERVICE_URL=http://localhost:8082
ENV SERVER_PORT=8083
EXPOSE 8083
ENTRYPOINT ["./presence-service"]

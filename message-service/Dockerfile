FROM golang:1.23.4 as builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o message-service ./cmd/message

FROM ubuntu:24.04
WORKDIR /app
COPY --from=builder /app/message-service .
COPY migrations ./migrations
ENV DATABASE_PATH=/app/messages.db
ENV SERVER_PORT=8081
EXPOSE 8081
ENTRYPOINT ["./message-service"]

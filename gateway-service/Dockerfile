FROM golang:1.23.4 as builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o gateway-service ./cmd/gateway

FROM ubuntu:24.04
WORKDIR /app
COPY --from=builder /app/gateway-service .
ENV AUTH_SERVICE_URL=http://localhost:8082
ENV MESSAGE_SERVICE_URL=http://localhost:8081
ENV SERVER_PORT=8080
EXPOSE 8080
ENTRYPOINT ["./gateway-service"]

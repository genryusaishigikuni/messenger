FROM golang:1.20-alpine as builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o auth-service ./cmd/auth

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/auth-service .
COPY migrations ./migrations
ENV DATABASE_PATH=/app/auth.db
ENV JWT_SECRET=changeme
EXPOSE 8082
ENTRYPOINT ["./auth-service"]

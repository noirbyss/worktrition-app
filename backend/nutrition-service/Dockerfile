FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY gen ./gen
COPY backend/nutrition-service/go.mod backend/nutrition-service/go.sum ./backend/nutrition-service/

WORKDIR /app/backend/nutrition-service

RUN go mod download

COPY backend/nutrition-service/cmd ./cmd
COPY backend/nutrition-service/internal ./internal
COPY backend/nutrition-service/migrations ./migrations

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/nutrition-service ./cmd

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /bin/nutrition-service /usr/local/bin/nutrition-service

EXPOSE 50052

CMD ["nutrition-service"]

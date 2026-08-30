FROM --platform=$BUILDPLATFORM golang:1.25.3-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY gen ./gen
COPY backend/user-service/go.mod backend/user-service/go.sum ./backend/user-service/

WORKDIR /app/backend/user-service

RUN go mod download

COPY backend/user-service/cmd ./cmd
COPY backend/user-service/internal ./internal
COPY backend/user-service/migrations ./migrations

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -o /bin/user-service ./cmd/user-service

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /bin/user-service /usr/local/bin/user-service

EXPOSE 50051

CMD ["user-service"]

FROM --platform=$BUILDPLATFORM golang:1.25.3-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY gen ./gen
COPY backend/gamification-service/go.mod backend/gamification-service/go.sum ./backend/gamification-service/

WORKDIR /app/backend/gamification-service

RUN go mod download

COPY backend/gamification-service/cmd ./cmd
COPY backend/gamification-service/internal ./internal
COPY backend/gamification-service/migrations ./migrations

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -o /bin/gamification-service ./cmd

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /bin/gamification-service /usr/local/bin/gamification-service

EXPOSE 50055

CMD ["gamification-service"]

FROM --platform=$BUILDPLATFORM golang:1.25.3-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY gen ./gen
COPY backend/workout-service/go.mod backend/workout-service/go.sum ./backend/workout-service/

WORKDIR /app/backend/workout-service

RUN go mod download

COPY backend/workout-service/cmd ./cmd
COPY backend/workout-service/internal ./internal

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -o /bin/workout-service ./cmd

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /bin/workout-service /usr/local/bin/workout-service

EXPOSE 50054

CMD ["workout-service"]

FROM --platform=$BUILDPLATFORM golang:1.25.3-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY gen ./gen
COPY backend/ai-service/go.mod backend/ai-service/go.sum ./backend/ai-service/

WORKDIR /app/backend/ai-service

RUN go mod download

COPY backend/ai-service/cmd ./cmd
COPY backend/ai-service/internal ./internal

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -o /bin/ai-service ./cmd/main.go

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /bin/ai-service /usr/local/bin/ai-service

EXPOSE 50056

CMD ["ai-service"]

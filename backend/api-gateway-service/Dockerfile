FROM --platform=$BUILDPLATFORM golang:1.25.3-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY gen ./gen
COPY backend/api-gateway-service/go.mod backend/api-gateway-service/go.sum ./backend/api-gateway-service/

WORKDIR /app/backend/api-gateway-service

RUN go mod download

COPY backend/api-gateway-service/cmd ./cmd
COPY backend/api-gateway-service/internal ./internal

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -o /bin/api-gateway-service ./cmd/api-gateway-service

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /bin/api-gateway-service /usr/local/bin/api-gateway-service

EXPOSE 8080

CMD ["api-gateway-service"]

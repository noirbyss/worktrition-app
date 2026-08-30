.PHONY: build run test docker-up docker-down docker-logs

build:
	go build -o bin/gamification-service ./cmd

run: build
	./bin/gamification-service

test:
	go test ./...

docker-up:
	docker compose up --build

docker-down:
	docker compose down -v

docker-logs:
	docker compose logs -f gamification-service

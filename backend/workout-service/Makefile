.PHONY: build run test docker-up docker-down docker-logs migrate-up migrate-down migrate-create

build:
	go build -o bin/workout-service ./cmd

run: build
	./bin/workout-service

test:
	go test ./...

docker-up:
	docker compose up --build

docker-down:
	docker compose down -v

docker-logs:
	docker compose logs -f workout-service

migrate-up:
	migrate -database "$$DATABASE_URL" -path migrations up

migrate-down:
	migrate -database "$$DATABASE_URL" -path migrations down 1

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

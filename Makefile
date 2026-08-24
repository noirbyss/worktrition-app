.PHONY: proto generate lint format

tools:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

generate:
	buf generate
	cd gen && go mod tidy

lint:
	buf lint

format:
	buf format -w
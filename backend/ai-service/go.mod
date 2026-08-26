module ai-service

go 1.25.3

replace github.com/noirbyss/worktrition-app/gen => ../../gen

require (
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/k0kubun/pp/v3 v3.5.2
	github.com/noirbyss/worktrition-app/gen v0.0.0-20260824195145-daf7f694f975
	github.com/sashabaranov/go-openai v1.42.0
	go.uber.org/zap v1.28.0
)

require (
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/grpc v1.83.1 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

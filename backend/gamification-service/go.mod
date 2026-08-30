module gamification-service

go 1.25.3

require (
	github.com/jackc/pgx/v5 v5.5.5
	github.com/noirbyss/worktrition-app/gen v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.83.1
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/crypto v0.51.0 // indirect
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace github.com/noirbyss/worktrition-app/gen => ../../gen

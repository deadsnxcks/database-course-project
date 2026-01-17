include .env
export

migrate-up:
	migrate -path migrations -database ${DB_CONN} up

migrate-down:
	migrate -path migrations -database ${DB_CONN} down

gen-proto:
	protoc \
		-I protos/proto \
		protos/proto/dbcp/dbcp.proto \
		--go_out=protos/gen/go \
		--go_opt=paths=source_relative \
		--go-grpc_out=protos/gen/go \
		--go-grpc_opt=paths=source_relative

run:
	go run cmd/dbcp/main.go
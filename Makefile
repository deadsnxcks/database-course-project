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

gen-vessel-proto:
	protoc \
		-I protos/proto \
		protos/proto/vessel/vessel.proto \
		--go_out=protos/gen/go \
		--go_opt=paths=source_relative \
		--go-grpc_out=protos/gen/go \
		--go-grpc_opt=paths=source_relative

gen-cargotype-proto:
	protoc \
		-I protos/proto \
		protos/proto/cargotype/cargotype.proto \
		--go_out=protos/gen/go \
		--go_opt=paths=source_relative \
		--go-grpc_out=protos/gen/go \
		--go-grpc_opt=paths=source_relative

gen-cargo-proto:
	protoc \
		-I protos/proto \
		protos/proto/cargo/cargo.proto \
		--go_out=protos/gen/go \
		--go_opt=paths=source_relative \
		--go-grpc_out=protos/gen/go \
		--go-grpc_opt=paths=source_relative

gen-operation-proto:
	protoc \
		-I protos/proto \
		protos/proto/operation/operation.proto \
		--go_out=protos/gen/go \
		--go_opt=paths=source_relative \
		--go-grpc_out=protos/gen/go \
		--go-grpc_opt=paths=source_relative

gen-storageloc-proto:
	protoc \
		-I protos/proto \
		protos/proto/storageloc/storageloc.proto \
		--go_out=protos/gen/go \
		--go_opt=paths=source_relative \
		--go-grpc_out=protos/gen/go \
		--go-grpc_opt=paths=source_relative

gen-opercargo-proto:
	protoc \
		-I protos/proto \
		protos/proto/opercargo/opercargo.proto \
		--go_out=protos/gen/go \
		--go_opt=paths=source_relative \
		--go-grpc_out=protos/gen/go \
		--go-grpc_opt=paths=source_relative

run:
	go run cmd/dbcp/main.go
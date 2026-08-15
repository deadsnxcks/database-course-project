.PHONY: up down down-v logs ps run-server run-gateway build \
        migrate-up migrate-down gen-proto \
        gen-report-proto gen-vessel-proto gen-cargotype-proto gen-cargo-proto \
        gen-operation-proto gen-storageloc-proto gen-opercargo-proto

-include .env
export

PROTO_DIR := protos/proto
GEN_DIR   := protos/gen/go

up:
	docker compose up -d --build

down:
	docker compose down

down-v:
	docker compose down -v

logs:
	docker compose logs -f

ps:
	docker compose ps

run-server:
	go run ./server/cmd/dbcp

run-gateway:
	go run ./gateway/cmd/dbcp-api-gateway

build:
	go build -o bin/server ./server/cmd/dbcp
	go build -o bin/gateway ./gateway/cmd/dbcp-api-gateway

migrate-up:
	migrate -path migrations -database "$(DB_CONN)" up

migrate-down:
	migrate -path migrations -database "$(DB_CONN)" down

gen-proto: gen-vessel-proto gen-cargotype-proto gen-cargo-proto \
           gen-operation-proto gen-storageloc-proto gen-opercargo-proto \
           gen-report-proto

gen-%-proto:
	protoc \
		-I $(PROTO_DIR) \
		$(PROTO_DIR)/$*/$*.proto \
		--go_out=$(GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) \
		--go-grpc_opt=paths=source_relative
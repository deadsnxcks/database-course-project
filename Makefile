include .env
export

migrate-up:
	migrate -path migrations -database ${DB_CONN} up

migrate-down:
	migrate -path migrations -database ${DB_CONN} down

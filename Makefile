GO111MODULE=on
CONFIG_PATH=config/config.json

.PHONY: up
up:
	docker-compose up -d db
	while true ; do docker-compose exec db pg_isready && break ; sleep 0.1; done
	make generate_swagger_docs
	go run main.go --config=$(CONFIG_PATH)

# Apply new migrations (if exist)
.PHONY: migrate
migrate:
	docker-compose up -d db
	while true ; do docker-compose exec db pg_isready && break ; sleep 0.1; done
	go run main.go --config=$(CONFIG_PATH) migrate

# drop old and create a new database
.PHONY: recreate_db
recreate_db:
	docker-compose rm -sf db
	docker-compose up -d db
	while true ; do docker-compose exec db pg_isready && break ; sleep 0.1; done
	make migrate

# console access to database
.PHONY: psql
psql:
	docker-compose exec db psql -U crud_web_service crud_web_service

# generates swagger documentation
.PHONY: generate_swagger_docs
generate_swagger_docs:
	go get github.com/swaggo/swag/cmd/swag
	swag init



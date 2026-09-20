# Load env file
DB_USER=clubdb
DB_PASS=clubdb_password
DB_HOST=localhost
DB_PORT=5435
DB_NAME=clubdb
DB_SSL_MODE=disable

MIGRATE=migrate
MIGRATIONS_DIR=database/migrations
DB_DSN=postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)
BINARY=go-base-service

# install dev tools
tool:
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${shell go env GOPATH}/bin v2.1.6

## Run migrations up
migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" up

## Rollback one migration
migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" down 1

migrate-version:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" version

migrate-force:
	@test -n "$(VERSION)" || (echo "VERSION is required. Usage: make migrate-force VERSION=1" && exit 1)
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" force $(VERSION)

migrate-create:
	@test -n "$(NAME)" || (echo "NAME is required. Usage: make migrate-create NAME=create_users_table" && exit 1)
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

up:
	docker compose up -d

down:
	docker compose down

run:
	go run cmd/main.go

test:
	go test ./...

test-race:
	go test -race ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

lint:
	golangci-lint run

vet:
	go vet ./...

mod:
	@go mod tidy
	@go mod vendor

gen:
	@swag init -g internal/framework/route/route.go -o internal/docs  --exclude pkg,db,deployment,scripts,vendor
	protoc --proto_path=proto --go_out=internal/infrastructure/grpc/generated --go_opt=paths=source_relative --go-grpc_out=internal/infrastructure/grpc/generated --go-grpc_opt=paths=source_relative proto/user/user.proto

build:
	@go build -tags=jsoniter -o $(BINARY) cmd/main.go

buildx:
	docker buildx create --use

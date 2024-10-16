set export
set dotenv-load

PLATFORM := if os() == "macos" { "darwin" } else { os() }

lint:
	golangci-lint run --fix

test:
	go test -covermode=count -coverprofile=count.out -v ./...

mock:
	@rm -rf mocks
	mockery --all

download-golang-migrate-binary:
	if [ ! -f ./migrate.{{PLATFORM}}-amd64 ] ; \
	then \
		curl -sfL https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.{{PLATFORM}}-amd64.tar.gz | tar -xvz; \
	fi;

db-migrate: download-golang-migrate-binary
	./migrate.{{PLATFORM}}-amd64 -source file://db/migrations -database $POSTMAND_DATABASE_URL up

db-test-migrate: download-golang-migrate-binary
	./migrate.{{PLATFORM}}-amd64 -source file://db/migrations -database $POSTMAND_TEST_DATABASE_URL up

run-server:
	go run cmd/postmand/main.go server

run-worker:
	go run cmd/postmand/main.go worker

swag-init:
	swag init -g cmd/postmand/main.go --parseDependency --parseDepth 2

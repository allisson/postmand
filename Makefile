PLATFORM := $(shell uname | tr A-Z a-z)

lint:
	if [ ! -f ./bin/golangci-lint ] ; \
	then \
		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.38.0; \
	fi;
	./bin/golangci-lint run

test:
	go test -covermode=count -coverprofile=count.out -v ./...

mock:
	@rm -rf mocks
	mockery --all

download-golang-migrate-binary:
	if [ ! -f ./migrate.$(PLATFORM)-amd64 ] ; \
	then \
		curl -sfL https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.$(PLATFORM)-amd64.tar.gz | tar -xvz; \
	fi;

db-migrate: download-golang-migrate-binary
	./migrate.$(PLATFORM)-amd64 -source file://db/migrations -database ${POSTMAND_DATABASE_URL} up

db-test-migrate: download-golang-migrate-binary
	./migrate.$(PLATFORM)-amd64 -source file://db/migrations -database ${POSTMAND_TEST_DATABASE_URL} up

run-server:
	go run cmd/postmand/main.go server

run-worker:
	go run cmd/postmand/main.go worker

swag-init:
	swag init -g cmd/postmand/main.go --parseDependency

.PHONY: lint test mock download-golang-migrate-binary db-migrate db-test-migrate run-server run-worker swag-init

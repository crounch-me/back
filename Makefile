DATABASE_PORT := 5432

VERSION := $(shell head -n 1 VERSION)
SEABOLT_DIR := $(shell pkg-config --variable=libdir seabolt17)

.PHONY: build
build:
	@echo "+ $@"
	go build -o $(GOPATH)/bin/crounch-back

.PHONY: acceptance-test
acceptance-test:
	@echo "+ $@"
	cd acceptance; godog

.PHONY: run
run: run-dependencies sleep run-app

.PHONY: sleep
sleep:
	sleep 3

.PHONY: run-app
run-app:
	@echo "+ $@"
	go run -ldflags "-r $(SEABOLT_DIR)" main.go 

.PHONY: run-dependencies
run-dependencies:
	@echo "+ $@"
	@docker-compose -p crounch-back -f containers/docker-compose.dependencies.yml down || true;
	@docker-compose -p crounch-back -f containers/docker-compose.dependencies.yml pull;
	@docker-compose -p crounch-back -f containers/docker-compose.dependencies.yml up -d --build

.PHONY: test
test:
	@echo "+ $@"
	go test -v -ldflags "-r $(SEABOLT_DIR)" ./handler 

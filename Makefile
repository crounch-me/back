DATABASE_PORT := 5432

APP_NAME := crounch-back
BUILDER_IMAGE_NAME := $(APP_NAME)-builder
TEST_IMAGE_NAME := $(APP_NAME)-test
DOCKER_USER := sehsyha

TAG_FLAG :=
ifneq ($(TAG),)
	TAG_FLAG:=--tags=$(TAG)
endif

.PHONY: build
build:
	@echo "+ $@"
	go build -o $(GOPATH)/bin/$(APP_NAME)

.PHONY: build-image
build-image: build-builder-image
	@echo "+ $@"
	docker build -f containers/Dockerfile -t $(APP_NAME) --build-arg BUILDER_IMAGE=$(BUILDER_IMAGE_NAME) .

.PHONY: build-builder-image
build-builder-image:
	@echo "+ $@"
	docker build -t $(BUILDER_IMAGE_NAME) -f containers/Dockerfile.builder .

.PHONY: acceptance-test
acceptance-test:
	@echo "+ $@"
	cd acceptance; godog $(TAG_FLAG)

.PHONY: acceptance-test-ci
acceptance-test-ci: run-image-ci
	@echo "+ $@"
	docker rm $(APP_NAME)-acceptance-test || true
	docker run --net='host' --name $(APP_NAME)-acceptance-test $(BUILDER_IMAGE_NAME) /bin/sh -c "make acceptance-test"

.PHONY: run
run: run-dependencies sleep run-app

.PHONY: run-app
run-app:
	@echo "+ $@"
	swag init -g router/router.go
	go run main.go serve --db-connection-uri postgresql://postgres:secretpassword@localhost/postgres?sslmode=disable --db-schema public

.PHONY: run-dependencies
run-dependencies:
	@echo "+ $@"
	@docker-compose -p $(APP_NAME) -f containers/docker-compose.dependencies.yml down || true;
	@docker-compose -p $(APP_NAME) -f containers/docker-compose.dependencies.yml pull;
	@docker-compose -p $(APP_NAME) -f containers/docker-compose.dependencies.yml up -d --build

.PHONY: run-image
run-image: build-builder-image run-image-ci

.PHONY: run-image-ci
run-image-ci:
	BUILDER_IMAGE_NAME=$(BUILDER_IMAGE_NAME) APP_NAME=$(APP_NAME) docker-compose -p $(APP_NAME) -f containers/docker-compose.yml up -d --build

.PHONY: test
test:
	@echo "+ $@"
	go test -v ./...

.PHONY: test-ci
test-ci: build-builder-image
	@echo "+ $@"
	docker run --name $(APP_NAME)-test $(BUILDER_IMAGE_NAME) /bin/sh -c 'make cover-ci'
	WORKDIR=$(shell docker inspect --format "{{.Config.WorkingDir}}" $(BUILDER_IMAGE_NAME)); \
		docker cp $(APP_NAME)-test:$$WORKDIR/profile.cov profile.cov

.PHONY: cover-ci
cover-ci:
	@echo "+ $@"
	go test -v -covermode=count -coverprofile=profile.cov ./handler

.PHONY: sleep
sleep:
	sleep 3

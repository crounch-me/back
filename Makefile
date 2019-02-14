DATABASE_PORT := 5432

VERSION := $(shell head -n 1 VERSION)
SEABOLT_DIR := $(shell pkg-config --variable=libdir seabolt17)

APP_NAME := crounch-back
BUILDER_IMAGE_NAME := $(APP_NAME)-builder:$(VERSION)
TEST_IMAGE_NAME := $(APP_NAME)-test-$(VERSION)
DOCKER_USER := sehsyha

.PHONY: bump-version
bump-version:
	@echo "+ $@"
	git fetch --tags
	echo '{"version": "$(VERSION)"}' > ./package.json
	npm i -g standard-version@4.2.0
	standard-version --skip.commit true --skip.tag true
	NEW_VERSION=`jq -r ".version" package.json`; \
		echo $$NEW_VERSION > VERSION; \
		git add CHANGELOG.md; \
		git add VERSION; \
		git commit -m "build: bump to version $$NEW_VERSION [skip ci]"; \
		git push origin master; \
		git tag $$NEW_VERSION; \
		git push --tags origin master; \
		rm package.json

.PHONY: build
build:
	@echo "+ $@"
	go build -o $(GOPATH)/bin/$(APP_NAME)

.PHONY: build-image
build-image: build-builder-image
	@echo "+ $@"
	docker build -f containers/Dockerfile -t $(APP_NAME):$(VERSION) --build-arg BUILDER_IMAGE=$(BUILDER_IMAGE_NAME) .
	docker tag $(APP_NAME):$(VERSION) $(DOCKER_USER)/$(APP_NAME):$(VERSION)

.PHONY: build-builder-image
build-builder-image:
	@echo "+ $@"
	docker build -t $(BUILDER_IMAGE_NAME) -f containers/Dockerfile.builder .

.PHONY: acceptance-test
acceptance-test:
	@echo "+ $@"
	cd acceptance; godog

.PHONY: acceptance-test-ci
acceptance-test-ci: build-builder-image run-image
	@echo "+ $@"
	docker rm $(APP_NAME)-acceptance-test || true
	docker run --net='host' --name $(APP_NAME)-acceptance-test $(BUILDER_IMAGE_NAME) /bin/sh -c "make acceptance-test"

.PHONY: run
run: run-dependencies sleep run-app

.PHONY: sleep
sleep:
	sleep 3

.PHONY: run-app
run-app:
	@echo "+ $@"
	go run -ldflags "-r $(SEABOLT_DIR)" main.go serve

.PHONY: run-dependencies
run-dependencies:
	@echo "+ $@"
	@docker-compose -p $(APP_NAME) -f containers/docker-compose.dependencies.yml down || true;
	@docker-compose -p $(APP_NAME) -f containers/docker-compose.dependencies.yml pull;
	@docker-compose -p $(APP_NAME) -f containers/docker-compose.dependencies.yml up -d --build

.PHONY: run-image
run-image: build-builder-image
	BUILDER_IMAGE_NAME=$(BUILDER_IMAGE_NAME) APP_NAME=$(APP_NAME) docker-compose -p $(APP_NAME) -f containers/docker-compose.yml up -d --build

.PHONY: test
test:
	@echo "+ $@"
	go test -v -ldflags "-r $(SEABOLT_DIR)"  $(shell go list ./... | grep -v vendor | grep -v acceptance)

.PHONY: cover-ci 
cover-ci:
	@echo "+ $@"
	go test -v -covermode=count -coverprofile=profile.cov ./handler

.PHONY: test-ci
test-ci: build-builder-image
	@echo "+ $@"
	docker rm $(APP_NAME)-test || true
	docker run --name $(APP_NAME)-test $(BUILDER_IMAGE_NAME) /bin/sh -c 'make cover-ci'
	WORKDIR=$(shell docker inspect --format "{{.Config.WorkingDir}}" $(BUILDER_IMAGE_NAME)); \
		docker cp $(APP_NAME)-test:$$WORKDIR/profile.cov profile.cov
DATABASE_PORT := 5432

VERSION := $(shell head -n 1 VERSION)

APP_NAME := crounch-back
BUILDER_IMAGE_NAME := $(APP_NAME)-builder:$(VERSION)
TEST_IMAGE_NAME := $(APP_NAME)-test-$(VERSION)
DOCKER_USER := sehsyha

TAG_FLAG :=
ifneq ($(TAG),)
	TAG_FLAG:=--tags=$(TAG)
endif

.PHONY: bump-version
bump-version:
	@echo "+ $@"
	git checkout master
	git fetch --tags
	echo '{"version": "$(VERSION)"}' > ./package.json
	npm i -g standard-version@4.2.0
	standard-version --skip.commit true --skip.tag true
	NEW_VERSION=`jq -r ".version" package.json`; \
		echo $$NEW_VERSION > VERSION; \
		git add CHANGELOG.md; \
		git add VERSION; \
		git commit -m "build: bump to version $$NEW_VERSION [skip ci]"; \
		git tag $$NEW_VERSION; \
		git remote rm origin; \
		git remote add origin https://$(DOCKER_USER):$(GH_TOKEN)@github.com/Sehsyha/crounch-back.git; \
		git push origin master; \
		git push --tags

.PHONY: build
build:
	@echo "+ $@"
	go build -o $(GOPATH)/bin/$(APP_NAME)

.PHONY: build-image
build-image: build-builder-image
	@echo "+ $@"
	docker build -f containers/Dockerfile -t $(APP_NAME):$(VERSION) --build-arg BUILDER_IMAGE=$(BUILDER_IMAGE_NAME) .
	docker tag $(APP_NAME):$(VERSION) $(DOCKER_USER)/$(APP_NAME):$(VERSION)
	docker tag $(APP_NAME):$(VERSION) $(DOCKER_USER)/$(APP_NAME):latest

.PHONY: build-builder-image
build-builder-image:
	@echo "+ $@"
	docker build -t $(BUILDER_IMAGE_NAME) -f containers/Dockerfile.builder .

.PHONY: publish-image
publish-image:
	@echo "+ $@"
	docker login -u $(DOCKER_USER) -p $(DOCKER_PASSWORD)
	docker push $(DOCKER_USER)/crounch-back:$(VERSION)
	docker push $(DOCKER_USER)/crounch-back:latest

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
	go test -v $(shell go list ./... | grep -v vendor | grep -v acceptance)

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

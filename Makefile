VERSION := $(shell jq -r ".version" package.json)

APP_NAME := crounch-front
BUILDER_IMAGE_NAME := $(APP_NAME)-builder:$(VERSION)
TEST_IMAGE_NAME := $(APP_NAME)-test-$(VERSION)
DOCKER_USER := sehsyha

.PHONY: bump-version
bump-version:
	@echo "+ $@"
	git checkout master
	git fetch --tags
	npm i -g standard-version@4.2.0
	standard-version --skip.commit true --skip.tag true
	NEW_VERSION=`jq -r ".version" package.json`; \
		git add CHANGELOG.md; \
		git add package.json; \
		git commit -m "build: bump to version $$NEW_VERSION [skip ci]"; \
		git tag $$NEW_VERSION; \
		git remote rm origin; \
		git remote add origin https://$(DOCKER_USER):$(GH_TOKEN)@github.com/Sehsyha/crounch-back.git; \
		git push origin master; \
		git push --tags

.PHONY: build-image
build-image:
	@echo "+ $@"
	docker build -f containers/Dockerfile -t $(APP_NAME):$(VERSION) .
	docker tag $(APP_NAME):$(VERSION) $(DOCKER_USER)/$(APP_NAME):$(VERSION)

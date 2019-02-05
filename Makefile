.PHONY: build
build:
	go build

.PHONY: clean
clean:
	rm ./crounch-back

.PHONY: run
run:
	go run main.go

.PHONY: test
test:
	go test -v ./handler

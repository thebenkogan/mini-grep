.PHONY: build

build:
	go build -o build/main ./...

test:
	go test ./...

.PHONY: build

build:
	go build -o build/main ./cmd/main.go

test:
	go test ./...

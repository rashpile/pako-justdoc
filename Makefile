BINARY_NAME=justdoc
BUILD_DIR=bin
PORT?=6001

.PHONY: build run test lint clean docker

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/justdoc

run: build
	PORT=$(PORT) ./$(BUILD_DIR)/$(BINARY_NAME)

test:
	go test -v -race -cover ./...

clean:
	rm -rf $(BUILD_DIR) *.db

lint:
	golangci-lint run

docker:
	docker build -t justdoc:latest .

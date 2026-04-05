.PHONY: run build proto tidy format install-tools clean

# Define project variables
PROJECT_NAME := go-microservice
PKG_LIST := $(shell go list ./...)

install-tools:
	@echo "Installing tools..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

proto:
	@echo "Generating protobuf definitions..."
	# Run this locally if protoc is installed, or adapt to use a Docker container:
	# docker run -v $(pwd):/workspace -w /workspace rvolosatovs/protoc:latest ...
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       api/proto/v1/user.proto

tidy:
	@echo "Tidying go mod..."
	go mod tidy -e

build:
	@echo "Building application..."
	go build -o bin/$(PROJECT_NAME) cmd/server/main.go

run:
	@echo "Running application..."
	go run cmd/server/main.go

format:
	@echo "Formatting code..."
	go fmt ./...

clean:
	@echo "Cleaning up..."
	rm -rf bin/

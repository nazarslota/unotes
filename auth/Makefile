PROJECT_NAME := auth

BIN_DIR := bin
PROTO_DIR := api/proto
SWAGGER_DIR := api/swagger

ifeq ($(OS), Windows_NT)
    $(Error OS not supported by this Makefile.)
else
    SHELL := bash
    SHELL_VERSION := $(shell echo $$BASH_VERSION)

    UNAME := $(shell uname -s)
    VERSION_AND_ARCH := $(shell uname -rm)

    ifeq ($(UNAME), Darwin)
    	OS := MacOS ${VERSION_AND_ARCH}
    else ifeq ($(UNAME), Linux)
    	OS := Linux ${VERSION_AND_ARCH}
    else
    	$(Error OS not supported by this Makefile.)
    endif

	GO_PACKAGE := $(shell head -1 go.mod | awk '{print $$2}')
    HELP_COMMAND := grep -E '^[a-zA-Z_-]+:.*?\#\# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?\#\# "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

    RM_F_COMMAND = rm -f
    RM_RF_COMMAND = ${RM_F_COMMAND} -r
endif

.DEFAULT_GOAL := help
.PHONY: clean build test swagger proto about help


build: clean test ## Generate build.
	go build -o ./${BIN_DIR}/ ./cmd/...

clean: ## Clean generated files.
	${RM_RF_COMMAND} ./${BIN_DIR}/
	${RM_F_COMMAND} ./${PROTO_DIR}/*.pb.go
	${RM_F_COMMAND} ./${SWAGGER_DIR}/*

test: dependencies proto swagger ## Launch tests.
	go test ./...

dependencies: ## Download all golang dependencies.
	go mod download

proto: ## Generate protocol buffer.
	protoc -I${PROTO_DIR} --go_out=. --go_opt=module=${GO_PACKAGE} --go-grpc_out=. --go-grpc_opt=module=${GO_PACKAGE} ${PROTO_DIR}/*.proto

swagger: ## Generate swagger doc.
	swag init -q -o ${SWAGGER_DIR} -g ./internal/handler/rest/handler.go

about: ## Display info related to the build.
	@echo "OS:              ${OS}"
	@echo "Shell:           ${SHELL} ${SHELL_VERSION}"
	@echo "Protoc version:  $(shell protoc --version)"
	@echo "Go version:      $(shell go version)"
	@echo "Go package:      ${GO_PACKAGE}"
	@echo "OpenSSL version: $(shell openssl version)"

help: ## Show this help.
	@${HELP_COMMAND}

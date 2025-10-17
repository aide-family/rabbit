GOHOSTOS:=$(shell go env GOHOSTOS)
VERSION=$(shell git describe --tags --always)
BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find proto/rabbit/api -name *.proto")
	# Use mkdir -p equivalent for Windows
	MKDIR=mkdir
	RM=del /f /q
else
	API_PROTO_FILES=$(shell find proto/rabbit/api -name *.proto)
	MKDIR=mkdir -p
	RM=rm -f
endif

.PHONY: init
# initialize the moon environment
init:
	@echo "Initializing moon environment"
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.3
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/moon-monitor/stringer@latest
	go install github.com/protoc-gen/i18n-gen@latest

.PHONY: conf
# generate the conf files
conf: config
	@echo "Generating conf files"
	protoc --proto_path=./internal/conf \
           --proto_path=./proto/rabbit \
           --proto_path=./proto/third_party \
           --go_out=paths=source_relative:./internal/conf \
           --experimental_allow_proto3_optional \
           ./internal/conf/*.proto

.PHONY: config
# generate the config files
config:
	@echo "Generating config files"
	@if [ "$(GOHOSTOS)" = "windows" ]; then \
		$(Git_Bash) -c "rm -rf ./pkg/config"; \
		if [ ! -d "./pkg/config" ]; then $(MKDIR) ./pkg/config; fi \
	else \
		rm -rf ./pkg/config; \
		if [ ! -d "./pkg/config" ]; then $(MKDIR) ./pkg/config; fi \
	fi
	protoc --proto_path=./proto/rabbit/config \
	       --proto_path=./proto/third_party \
	       --go_out=paths=source_relative:./pkg/config \
	       --experimental_allow_proto3_optional \
	       ./proto/rabbit/config/*.proto

.PHONY: api
# generate the api files
api:
	@echo "Generating api files"
	@if [ "$(GOHOSTOS)" = "windows" ]; then \
		$(Git_Bash) -c "rm -rf ./pkg/api"; \
		if [ ! -d "./pkg/api" ]; then $(MKDIR) ./pkg/api; fi \
	else \
		rm -rf ./pkg/api; \
		if [ ! -d "./pkg/api" ]; then $(MKDIR) ./pkg/api; fi \
	fi
	protoc --proto_path=./proto/rabbit/api \
	       --proto_path=./proto/third_party \
 	       --go_out=paths=source_relative:./pkg/api \
 	       --go-http_out=paths=source_relative:./pkg/api \
 	       --go-grpc_out=paths=source_relative:./pkg/api \
	       --openapi_out=fq_schema_naming=true,default_response=false:./internal/server/swagger \
	       --experimental_allow_proto3_optional \
	       $(API_PROTO_FILES)

.PHONY: i18n
# i18n generate the i18n files
i18n:
	i18n-gen -O ./i18n/ -P ./proto/rabbit/**.proto -L en,zh -suffix Error

.PHONY: errors
# generate errors
errors:
	@echo "Generating errors"
	@if [ "$(GOHOSTOS)" = "windows" ]; then \
		$(Git_Bash) -c "rm -rf ./pkg/merr"; \
		if [ ! -d "./pkg/merr" ]; then $(MKDIR) ./pkg/merr; fi \
	else \
		rm -rf ./pkg/merr; \
		if [ ! -d "./pkg/merr" ]; then $(MKDIR) ./pkg/merr; fi \
	fi
	protoc --proto_path=./proto/rabbit/merr \
           --proto_path=./proto/third_party \
           --go_out=paths=source_relative:./pkg/merr \
           --go-errors_out=paths=source_relative:./pkg/merr \
           ./proto/rabbit/merr/*.proto
	make i18n

.PHONY: wire
# generate the wire files
wire:
	@echo "Generating wire files"
	wire ./...

.PHONY: build
# build the rabbit binary
build:
	@echo "Building rabbit"
	go build -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)" -o bin/rabbit main.go

.PHONY: dev
# run the rabbit binary in development mode
dev:
	@echo "Running rabbit in development mode"
	go run . run --swagger --metrics

.PHONY: test
# run the tests
test:
	@echo "Running tests"
	go test ./...

.PHONY: clean
# clean the binary
clean:
	@echo "Cleaning up"
	rm -rf bin

.PHONY: completion
# generate the completion
completion:
	@echo "Generating completion"
	go run . completion bash > completion.bash

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
GOHOSTOS:=$(shell go env GOHOSTOS)
VERSION=$(shell git describe --tags --always)
BUILD_TIME=$(shell date '+%Y-%m-%dT%H:%M:%SZ')
AUTHOR=$(shell git log -1 --format='%an')
AUTHOR_EMAIL=$(shell git log -1 --format='%ae')
REPO=$(shell git config remote.origin.url)

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
	go install golang.org/x/tools/gopls@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest

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
		$(Git_Bash) -c "rm -rf ./pkg/config/*.pb.go"; \
		if [ ! -d "./pkg/config" ]; then $(MKDIR) ./pkg/config; fi \
	else \
		rm -rf ./pkg/config/*.pb.go; \
		if [ ! -d "./pkg/config" ]; then $(MKDIR) ./pkg/config; fi \
	fi
	protoc --proto_path=./proto/rabbit/config \
	       --proto_path=./proto/rabbit \
	       --proto_path=./proto/third_party \
	       --go_out=paths=source_relative:./pkg/config \
	       --experimental_allow_proto3_optional \
	       ./proto/rabbit/config/*.proto

.PHONY: enum
# generate the enum files
enum:
	@echo "Generating enum files"
	@if [ "$(GOHOSTOS)" = "windows" ]; then \
		$(Git_Bash) -c "rm -rf ./pkg/enum/*.pb.go"; \
		if [ ! -d "./pkg/enum" ]; then $(MKDIR) ./pkg/enum; fi \
	else \
		rm -rf ./pkg/enum/*.pb.go; \
		if [ ! -d "./pkg/enum" ]; then $(MKDIR) ./pkg/enum; fi \
	fi
	protoc --proto_path=./proto/rabbit/enum \
	       --proto_path=./proto/third_party \
	       --go_out=paths=source_relative:./pkg/enum \
	       --experimental_allow_proto3_optional \
	       ./proto/rabbit/enum/*.proto

.PHONY: api
# generate the api files
api: enum
	@echo "Generating api files"
	@if [ "$(GOHOSTOS)" = "windows" ]; then \
		$(Git_Bash) -c "rm -rf ./pkg/api/*.pb.go"; \
		if [ ! -d "./pkg/api" ]; then $(MKDIR) ./pkg/api; fi \
	else \
		rm -rf ./pkg/api/*.pb.go; \
		if [ ! -d "./pkg/api" ]; then $(MKDIR) ./pkg/api; fi \
	fi
	protoc --proto_path=./proto/rabbit/api \
	       --proto_path=./proto/rabbit \
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
		$(Git_Bash) -c "rm -rf ./pkg/merr/*.pb.go"; \
		if [ ! -d "./pkg/merr" ]; then $(MKDIR) ./pkg/merr; fi \
	else \
		rm -rf ./pkg/merr/*.pb.go; \
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

.PHONY: vobj
# generate the vobj files
vobj:
	@echo "Generating vobj files"
	cd internal/biz/vobj && go generate .

.PHONY: gorm-gen
# generate the gorm files
gorm-gen:
	@echo "Generating gorm files"
	go run ./cmd/gorm gorm gen

.PHONY: gorm-migrate
# migrate the gorm files
gorm-migrate:
	@echo "Migrating gorm files"
	go run ./cmd/gorm gorm migrate

.PHONY: all
# generate all files
all: 
	@git log -1 --format='%B' > description.txt
	make clean errors api conf vobj gorm-gen wire

.PHONY: build
# build the rabbit binary
build: all
	@echo "Building rabbit"
	@echo "VERSION: $(VERSION)"
	@echo "BUILD_TIME: $(BUILD_TIME)"
	@echo "AUTHOR: $(AUTHOR)"
	@echo "AUTHOR_EMAIL: $(AUTHOR_EMAIL)"
	@git log -1 --format='%B' > description.txt
	go build -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.Author=$(AUTHOR) -X main.Email=$(AUTHOR_EMAIL) -X main.Repo=$(REPO)" -o bin/rabbit main.go

.PHONY: dev
# run the rabbit binary in development mode
dev:
	@echo "Running rabbit in development mode"
	go run . run

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
	rm -rf internal/biz/do/query
	rm -rf internal/biz/vobj/*__string.go
	rm -rf internal/conf/*.pb.go
	rm -rf pkg/api/*/*.pb.go
	rm -rf pkg/api/*/*.pb.gw.go
	rm -rf pkg/enum/*.pb.go
	rm -rf pkg/merr/*.pb.go
	rm -rf pkg/config/*.pb.go

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
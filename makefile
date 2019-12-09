.PHONY: all example-server fmt start run stop docker test clean

GOPATH:=$(shell go env GOPATH)
ROOT_DIR = $(CURDIR)
BIN_DIR = $(ROOT_DIR)/bin
TOOLS_DIR = $(ROOT_DIR)/tools

CURR_TIME = $(shell date "+%Y-%m-%d %H:%M:%S")
HOST_NAME = $(shell hostname)
GO_VERSION = $(shell go version)
VERSION=0.1.0
SERVER_NAME=example-server

LDFLAGS = "-X '$(SERVER_NAME)/pkg/version.Version=${VERSION}' -X '$(SERVER_NAME)/pkg/version.BuildTime=${CURR_TIME}' -X '$(SERVER_NAME)/pkg/version.BuildHost=${HOST_NAME}' -X '$(SERVER_NAME)/pkg/version.GOVersion=${GO_VERSION}'"

all: $(SERVER_NAME)

$(SERVER_NAME):
	CGO_ENABLED=0 GO111MODULE=on GOFLAGS=-mod=vendor go build -ldflags $(LDFLAGS) -o $(BIN_DIR)/$(SERVER_NAME) ./cmd/$(SERVER_NAME)

fmt:
	go fmt ./cmd/...
	go fmt ./pkg/...
	go fmt ./internal/...

run: $(SERVER_NAME)
	export EXAMPLE_SERVER_LOGPATH=/tmp \
	&& export EXAMPLE_SERVER_MYSQLHOST=127.0.0.1:3306 \
	&& export EXAMPLE_SERVER_MYSQLDB=test \
	&& export EXAMPLE_SERVER_MYSQLUSER=test \
	&& export EXAMPLE_SERVER_MYSQLPASSWORD=123456 \
	&& export EXAMPLE_SERVER_REDISMODE=standalone \
	&& export EXAMPLE_SERVER_REDISHOST=172.17.0.1:6379 \
	&& export EXAMPLE_SERVER_REDISDB=2 \
	&& $(BIN_DIR)/$(SERVER_NAME)

start: stop
	export EXAMPLE_SERVER_LOGPATH=/tmp \
	&& export EXAMPLE_SERVER_MYSQLHOST=127.0.0.1:3306 \
	&& export EXAMPLE_SERVER_MYSQLDB=test \
	&& export EXAMPLE_SERVER_MYSQLUSER=test \
	&& export EXAMPLE_SERVER_MYSQLPASSWORD=123456 \
	&& export EXAMPLE_SERVER_REDISMODE=standalone \
	&& export EXAMPLE_SERVER_REDISHOST=172.17.0.1:6379 \
	&& export EXAMPLE_SERVER_REDISDB=2 \
	&& nohup $(BIN_DIR)/$(SERVER_NAME)> /dev/null &

stop:
	$(TOOLS_DIR)/kill_app.sh $(SERVER_NAME)

install:
	cp -rf $(BIN_DIR)/* ~/bin/

docker:
	cp -rf ./build/Dockerfile ./
	docker build . -t test/$(SERVER_NAME):${VERSION}
	rm -rf ./Dockerfile

test:
	go test -v ./cmd/... ./internal/... ./pkg/...

clean:
	rm -rf $(BIN_DIR)

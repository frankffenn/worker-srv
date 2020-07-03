PROTOC_INCLUDES=-I$(GOPATH)/src -I./proto/

all: proto build
.PHONY: all

proto:
	protoc ${PROTOC_INCLUDES} --go_out=plugins=micro,paths=source_relative:./proto ./proto/worker.proto
.PHONY: proto

build:
	rm -f worker-srv
	go build $(GOFLAGS) -o worker-srv 
.PHONY: build
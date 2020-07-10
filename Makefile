PROTOC_INCLUDES=-I$(GOPATH)/src -I./proto/

all: proto build
.PHONY: all

proto:
	protoc ${PROTOC_INCLUDES} --micro_out=,paths=source_relative:./proto --go_out=,paths=source_relative:./proto ./proto/worker.proto
	cd registry/service;\
	protoc ${PROTOC_INCLUDES} --micro_out=,paths=source_relative:./proto --go_out=,paths=source_relative:./proto ./proto/registry.proto
.PHONY: proto

build:
	rm -f worker-srv
	go build $(GOFLAGS) -o worker-srv src/*.go
	rm -f registry-srv
	go build $(GOFLAGS) -o registry-srv registry/*.go
.PHONY: build
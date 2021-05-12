.PHONY: all
all: proto binaries

define build_binary
	go build -o $(1).bin ./cmd/$(1)
endef

.PHONY: binaries
binaries: client.bin server.bin

.PHONY: client.bin
client.bin:
	$(call build_binary,$(@:.bin=))

.PHONY: server.bin
server.bin:
	$(call build_binary,$(@:.bin=))

.PHONY: proto
proto:
	protoc --gofast_out=plugins=grpc:. --plugin protoc-gen-gofast=$(shell which protoc-gen-gofast) ./proto/feature.proto 

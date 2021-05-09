.PHONY: proto
proto:
	protoc --gofast_out=plugins=grpc:. --plugin protoc-gen-gofast=$(shell which protoc-gen-gofast) ./proto/feature.proto 

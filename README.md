# goff

Feature flags in Go, backed by a grpc service.

Surely not a bad idea.

## Quickstart

1. Build the server.

    `$ go build -o server.bin ./cmd/server`
1. Build the client.

    `$ go build -o client.bin ./cmd/client`
1. In one tab, run the server.

    `$ ./server.bin`

1. In a second tab, enable, disable, and remove features.

```
$ ./client.bin set foo true
$ ./client.bin set bar false
$ ./client.bin list-all
foo:true
bar:false
$ ./client.bin set foo false
$ ./client.bin delete bar
deleted feature bar:false
$ ./client.bin delete baz
no such feature baz
$ ./client.bin list-all
foo:false
```

## Development

1. [Install protoc](https://grpc.io/docs/protoc-installation/).

1. Install latest `protoc-gen-gofast`.

    `$ go install github.com/gogo/protobuf/protoc-gen-gofast@latest`

1. If changing anything in `./proto/*.proto`, re-generate the generated Go code.

    `$ make proto`

## "Roadmap"

See [TODO.md](./TODO.md).

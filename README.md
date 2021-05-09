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
$ ./client.bin set foo constant --enabled
$ ./client.bin set bar percentage_based -p10
$ ./client.bin set baz expression -e"1 > 0"
$ ./client.bin list
foo:true
bar:false
baz:false
$ ./client.bin list -j | jq '.'
{
  "bar": {
    "name": "bar",
    "type": "PERCENTAGE_BASED",
    "percentage": 10
  },
  "baz": {
    "name": "baz",
    "type": "EXPRESSION",
    "expression": "1 > 0"
  },
  "foo": {
    "name": "foo",
    "type": "CONSTANT",
    "enabled": true
  }
}
$ ./client.bin delete foo
deleted feature bar:true
$ ./client.bin delete qux
no such feature qux
$ ./client.bin list -j | jq '.'
{
  "bar": {
    "name": "bar",
    "type": "PERCENTAGE_BASED",
    "percentage": 10
  },
  "baz": {
    "name": "baz",
    "type": "EXPRESSION",
    "expression": "1 > 0"
  }
}
```

### Persistence

The feature server stores feature state in memory (for now). To get a modicum
of persistence, you can use the client to dump out a JSON representation of the
feature map, and then use that file when restarting the server.

For example:

```
prompt1> $ ./server.bin
prompt2> $ ./client.bin set foo constant --enabled
prompt2> $ ./client.bin set bar percentage_based -p10
prompt2> $ ./client.bin set baz expression -e"1 > 0"
prompt2> $ ./client.bin list -j > feature_flags.json
prompt1> $ ./server.bin --config feature_flags.json
prompt2> $ ./client.bin list -j
{
  "bar": {
    "name": "bar",
    "type": "PERCENTAGE_BASED",
    "percentage": 10
  },
  "baz": {
    "name": "baz",
    "type": "EXPRESSION",
    "expression": "1 > 0"
  },
  "foo": {
    "name": "foo",
    "type": "CONSTANT",
    "enabled": true
  }
}
```

## Development

1. [Install protoc](https://grpc.io/docs/protoc-installation/).

1. Install latest `protoc-gen-gofast`.

    `$ go install github.com/gogo/protobuf/protoc-gen-gofast@latest`

1. If changing anything in `./proto/*.proto`, re-generate the generated Go code.

    `$ make proto`

## "Roadmap"

See [TODO.md](./TODO.md).

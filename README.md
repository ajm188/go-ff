# go-ff

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

### In your code

```go
import "github.com/ajm188/go-ff/feature"

// Initialization code
func init() {
    if err := feature.InitFromFile("/path/to/feature_flags.json"); err != nil {
        panic(err)
    }
}

func SomeFunction() {
    enabled, err := feature.Get("name_of_my_feature", nil)
    if err != nil {
        // handle error
    }

    if enabled {
        DoNewThing()
    } else {
        DoOldThing()
    }
}

func SomeFunctionWithParameters(name string) {
    enabled, err := feature.Get("another_feature", map[string]interface{
        "name": name,
    })
    if err != nil {
        // handle error
    }

    if enabled {
        DoNewThing()
    } else {
        DoOldThing()
    }
}
```

For dynamic modification of feature flags, you would also want to run a gRPC
server that has the FeaturesServer service running, e.g.:

```go
import (
    "github.com/ajm188/go-ff/feature"
    "google.golang.org/grpc"
)

func run() {
    s := grpc.NewServer()
    feature.RegisterServer(s)

    // Register your other services on `s`.

    go s.Serve(lis)
}
```

And now you can use the client, as in the [quickstart](#Quickstart), to add,
modify, and remove features from your code at runtime.

For more advanced production uses, you will want to run the FeatureServer gRPC
service on a separate port to restrict access to administrators and operators
of your service (i.e. you should not have "feature flag admin" on the same port
as your actual application rpcs).

## Development

1. [Install protoc](https://grpc.io/docs/protoc-installation/).

1. Install latest `protoc-gen-gofast`.

    `$ go install github.com/gogo/protobuf/protoc-gen-gofast@latest`

1. If changing anything in `./proto/*.proto`, re-generate the generated Go code.

    `$ make proto`

## "Roadmap"

See [TODO.md](./TODO.md).

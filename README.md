# squzy_go

Squzy_go is a package, which allows you to monitor you golang applications networking. 
It can be easily used via GRPC, gin and http requests.

## Application setup

Before call the API methods, it is necessary to setup application with two parameters:

```go
import "github.com/squzy/squzy_go/core"
```

```go
///client - your http client
application, err := core.CreateApplication(client, &core.Options{
	Name:       "application name",
	Host:       "host adress",
})
```

The `AgentId` parameter in options is not used yet.

## GRPC integration

To use the squzy monitoring with GRPC, import:

```go
import "github.com/squzy/squzy_go/squzy_grpc"
```

Squzy monitoring allows you to use it on both client and server side.
You don't need to duplicate it, so, if you use it on server side, you don't need to use it on client.



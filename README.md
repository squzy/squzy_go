# squzy_go

Squzy_go is a package, which allows you to monitor you golang applications networking. 
It can be easily used via GRPC, gin and http requests.

## Application setup

Before calling the API methods, it is necessary to setup application with two parameters:

```go
import squzy_core "github.com/squzy/squzy_go/core"
```

```go
//client - your http client, can be nil
app, err := squzy_core.CreateApplication(client, &squzy_core.Options{
		ApiHost:         "your squzy api host",
		ApplicationName: "your applciation name",
		ApplicationHost: "your applciation host",
	})
```

The `AgentId` parameter in options is not used yet.

## GRPC integration

To use the squzy monitoring with GRPC, import:

```go
import squzy_grpc "github.com/squzy/squzy_go/integrations/grpc"
```

Squzy monitoring allows you to use it on both client and server side.
You don't need to duplicate it, so, if you use it on server side, you don't need to use it on client.

Squzy is working through interceptor interfaces provided by GRPC. 
There are two implementations: unary interceptor and stream interceptor.
To use the interceptor you need to define application as it was mentioned above.

Client side usage:

```go
conn, err := grpc.Dial(
    grpcUri, 
    squzy_grpc.NewClientUnaryInterceptor(application),
    squzy_grpc.NewClientStreamUnaryInterceptor(application)
)
```

Server side usage:

```go
server := grpc.NewServer(
    squzy_grpc.NewServerUnaryInterceptor(application),
    squzy_grpc.NewServerStreamInterceptor(application),
)
```

## Gin integration

To use the Squzy monitoring with GRPC, import:

```go
import squzy_gin "github.com/squzy/squzy_go/integrations/gin"
```

Then you need to add Squzy middleware in your gin.Engine:

```go
r := gin.New()
r.Use(squyz_gin.New(application))
```

## http integration

To use Squzy monitoring with http, import:

```go
import squzy_http "github.com/squzy/squzy_go/integrations/http"
```

The Squzy gor http working through the `http.RoundTripper`.
You can provide you basic `http.RoundTripper` as `parent` parameter, or set it as `nil`.

```go
client := &http.Client{
    Transport: NewRoundTripper(application, parent),
}
```
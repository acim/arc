# Skeleton for REST service including RenderJSON, ZapLogger and PromMetrics middlewares

## Included middlewares:

* RenderJSON - simplifies implementation of JSON REST API's
* ZapLogger - [chi](https://github.com/go-chi/chi) middleware for logging using [zap](https://github.com/uber-go/zap) logger
* PromMetrics - [chi](https://github.com/go-chi/chi) middleware providing [Prometheus](https://prometheus.io/) metrics to your HTTP server
  Tracks total number of requests and requests duration partitioned by status code, method and request URI

## RenderJSON handler example

```go
import (
    "github.com/acim/pkg/middleware"
)

func ExampleHandler(w http.ResponseWriter, r *http.Request) {
    res := middleware.ResponseFromContext(r.Context())
    payload := &struct{
        foo string
        bar string
    }{"example", "golang"}
    res.SetPayload(payload).SetStatusCode(http.StatusAccepted)
}
```

That's all, your response will be encoded as application/json.

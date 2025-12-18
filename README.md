# API Platform Gateway Extensions

This repository provides a collection of gateway policies for the WSO2 API Platform. These policies allow you to customize and control API behavior at the gateway level.

## Features

### Available Policies

#### 1. Add Header Policy
Adds HTTP headers to requests or responses flowing through the gateway.

**Use Cases:**
- Add custom headers for backend services
- Include correlation IDs for tracing
- Add metadata headers for downstream processing

**Example:**
```go
import "github.com/renuka-fernando/api-platform-gateway-extensions/pkg/policies"

// Create an add header policy for requests
policy := policies.NewAddHeaderPolicy("X-Custom-Header", "my-value", policies.RequestFlow)

// Apply to a request
ctx := &policies.PolicyContext{Request: req}
err := policy.Execute(ctx)
```

#### 2. Remove Header Policy
Removes specific HTTP headers from requests or responses.

**Use Cases:**
- Remove sensitive headers before forwarding to backend
- Strip internal headers before sending response to client
- Clean up unnecessary headers for security

**Example:**
```go
// Create a remove header policy for responses
policy := policies.NewRemoveHeaderPolicy("X-Internal-Token", policies.ResponseFlow)

// Apply to a response
ctx := &policies.PolicyContext{Response: resp}
err := policy.Execute(ctx)
```

#### 3. Interceptor Policy
Calls an external service for custom request/response processing and header manipulation.

**Use Cases:**
- Dynamic header manipulation based on external logic
- Integration with external validation or enrichment services
- Complex transformations requiring external computation

**Example:**
```go
// Create an interceptor policy
policy := policies.NewInterceptorPolicy("http://interceptor-service:8080/process", policies.RequestFlow)

// Apply to a request
ctx := &policies.PolicyContext{Request: req}
err := policy.Execute(ctx)
```

### Policy Executor

Chain multiple policies together for complex scenarios:

```go
executor := policies.NewExecutor()

// Add multiple policies
executor.AddPolicy(policies.NewAddHeaderPolicy("X-Request-ID", "12345", policies.RequestFlow))
executor.AddPolicy(policies.NewRemoveHeaderPolicy("X-Internal-Token", policies.RequestFlow))

// Execute all policies in order
ctx := &policies.PolicyContext{Request: req}
err := executor.Execute(ctx)
```

## Installation

```bash
go get github.com/renuka-fernando/api-platform-gateway-extensions
```

## Usage

### Basic Policy Usage

```go
package main

import (
    "net/http"
    "github.com/renuka-fernando/api-platform-gateway-extensions/pkg/policies"
)

func main() {
    // Create a policy
    addHeaderPolicy := policies.NewAddHeaderPolicy("X-API-Version", "v1", policies.RequestFlow)
    
    // Validate the policy
    if err := addHeaderPolicy.Validate(); err != nil {
        panic(err)
    }
    
    // Create a request
    req, _ := http.NewRequest("GET", "http://api.example.com/users", nil)
    
    // Apply the policy
    ctx := &policies.PolicyContext{Request: req}
    if err := addHeaderPolicy.Execute(ctx); err != nil {
        panic(err)
    }
    
    // The request now has the X-API-Version header set
}
```

### Chaining Multiple Policies

```go
func applyPolicies(req *http.Request) error {
    executor := policies.NewExecutor()
    
    // Add policies in desired order
    executor.AddPolicy(policies.NewAddHeaderPolicy("X-Request-ID", generateRequestID(), policies.RequestFlow))
    executor.AddPolicy(policies.NewAddHeaderPolicy("X-Gateway", "api-platform", policies.RequestFlow))
    executor.AddPolicy(policies.NewRemoveHeaderPolicy("Authorization", policies.RequestFlow))
    
    // Execute all policies
    ctx := &policies.PolicyContext{Request: req}
    return executor.Execute(ctx)
}
```

### Using with HTTP Middleware

```go
func policyMiddleware(policies *policies.Executor) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx := &policies.PolicyContext{Request: r}
            if err := policies.Execute(ctx); err != nil {
                http.Error(w, "Policy execution failed", http.StatusInternalServerError)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

## Policy Flows

Policies can be applied at different stages:

- **RequestFlow**: Applied before the request reaches the backend
- **ResponseFlow**: Applied before the response is returned to the client
- **FaultFlow**: Applied when an error occurs (future enhancement)

## Development

### Running Tests

```bash
go test ./pkg/policies/... -v
```

### Building

```bash
go build ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
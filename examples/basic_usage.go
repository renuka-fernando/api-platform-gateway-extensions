package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/renuka-fernando/api-platform-gateway-extensions/pkg/policies"
)

func main() {
	fmt.Println("API Platform Gateway Extensions - Basic Usage Example")
	fmt.Println("======================================================\n")

	// Example 1: Add Header Policy
	fmt.Println("Example 1: Add Header Policy")
	addHeaderExample()
	fmt.Println()

	// Example 2: Remove Header Policy
	fmt.Println("Example 2: Remove Header Policy")
	removeHeaderExample()
	fmt.Println()

	// Example 3: Chaining Multiple Policies
	fmt.Println("Example 3: Chaining Multiple Policies")
	chainingExample()
	fmt.Println()

	// Example 4: HTTP Middleware Integration
	fmt.Println("Example 4: HTTP Middleware Integration")
	middlewareExample()
}

func addHeaderExample() {
	// Create an add header policy
	policy := policies.NewAddHeaderPolicy("X-API-Version", "v1.0", policies.RequestFlow)

	// Create a sample request
	req := httptest.NewRequest("GET", "http://api.example.com/users", nil)
	fmt.Printf("Before: Headers = %v\n", req.Header)

	// Apply the policy
	ctx := &policies.PolicyContext{Request: req}
	if err := policy.Execute(ctx); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("After: Headers = %v\n", req.Header)
	fmt.Printf("X-API-Version = %s\n", req.Header.Get("X-API-Version"))
}

func removeHeaderExample() {
	// Create a remove header policy
	policy := policies.NewRemoveHeaderPolicy("Authorization", policies.RequestFlow)

	// Create a sample request with Authorization header
	req := httptest.NewRequest("GET", "http://api.example.com/users", nil)
	req.Header.Set("Authorization", "Bearer secret-token")
	req.Header.Set("Content-Type", "application/json")
	fmt.Printf("Before: Headers = %v\n", req.Header)

	// Apply the policy
	ctx := &policies.PolicyContext{Request: req}
	if err := policy.Execute(ctx); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("After: Headers = %v\n", req.Header)
	fmt.Printf("Authorization header removed: %v\n", req.Header.Get("Authorization") == "")
}

func chainingExample() {
	// Create an executor
	executor := policies.NewExecutor()

	// Add multiple policies
	executor.AddPolicy(policies.NewAddHeaderPolicy("X-Request-ID", "12345-67890", policies.RequestFlow))
	executor.AddPolicy(policies.NewAddHeaderPolicy("X-Gateway", "api-platform", policies.RequestFlow))
	executor.AddPolicy(policies.NewRemoveHeaderPolicy("X-Internal-Token", policies.RequestFlow))

	// Create a sample request
	req := httptest.NewRequest("GET", "http://api.example.com/users", nil)
	req.Header.Set("X-Internal-Token", "internal-secret")
	req.Header.Set("User-Agent", "TestClient/1.0")
	fmt.Printf("Before: Headers = %v\n", req.Header)

	// Execute all policies
	ctx := &policies.PolicyContext{Request: req}
	if err := executor.Execute(ctx); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("After: Headers = %v\n", req.Header)
	fmt.Println("Policies applied:")
	for _, p := range executor.GetPolicies() {
		fmt.Printf("  - %s\n", p.Name())
	}
}

func middlewareExample() {
	// Create policy executor
	executor := policies.NewExecutor()
	executor.AddPolicy(policies.NewAddHeaderPolicy("X-Gateway", "api-platform", policies.RequestFlow))

	// Create middleware
	middleware := policyMiddleware(executor)

	// Create a simple handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Handler received request with headers: %v\n", r.Header)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	})

	// Wrap handler with middleware
	wrappedHandler := middleware(handler)

	// Create test request
	req := httptest.NewRequest("GET", "http://api.example.com/test", nil)
	rr := httptest.NewRecorder()

	// Serve request
	wrappedHandler.ServeHTTP(rr, req)
	fmt.Printf("Response status: %d\n", rr.Code)
}

// policyMiddleware creates an HTTP middleware that applies policies
func policyMiddleware(executor *policies.Executor) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := &policies.PolicyContext{Request: r}
			if err := executor.Execute(ctx); err != nil {
				http.Error(w, "Policy execution failed", http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

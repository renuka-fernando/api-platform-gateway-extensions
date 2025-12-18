package policies

import "net/http"

// PolicyContext contains the context information for policy execution
type PolicyContext struct {
	Request  *http.Request
	Response *http.Response
	Headers  http.Header
	// Metadata for additional policy-specific information
	Metadata map[string]interface{}
}

// Policy defines the interface that all policies must implement
type Policy interface {
	// Name returns the name of the policy
	Name() string
	
	// Execute applies the policy to the given context
	Execute(ctx *PolicyContext) error
	
	// Validate validates the policy configuration
	Validate() error
}

// Flow represents the flow in which a policy can be applied
type Flow string

const (
	// RequestFlow is applied before the request reaches the backend
	RequestFlow Flow = "request"
	
	// ResponseFlow is applied before the response is returned to the client
	ResponseFlow Flow = "response"
	
	// FaultFlow is applied when an error occurs
	FaultFlow Flow = "fault"
)

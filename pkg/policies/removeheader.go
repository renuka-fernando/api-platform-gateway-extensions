package policies

import (
	"errors"
	"net/http"
)

// RemoveHeaderPolicy removes HTTP headers from request or response
type RemoveHeaderPolicy struct {
	// HeaderName is the name of the header to remove
	HeaderName string
	// Flow specifies whether to apply on request or response
	Flow Flow
}

// NewRemoveHeaderPolicy creates a new RemoveHeaderPolicy
func NewRemoveHeaderPolicy(headerName string, flow Flow) *RemoveHeaderPolicy {
	return &RemoveHeaderPolicy{
		HeaderName: headerName,
		Flow:       flow,
	}
}

// Name returns the name of the policy
func (p *RemoveHeaderPolicy) Name() string {
	return "RemoveHeader"
}

// Execute applies the remove header policy to the context
func (p *RemoveHeaderPolicy) Execute(ctx *PolicyContext) error {
	if ctx == nil {
		return errors.New("policy context cannot be nil")
	}

	var headers http.Header
	
	switch p.Flow {
	case RequestFlow:
		if ctx.Request != nil {
			headers = ctx.Request.Header
		}
	case ResponseFlow:
		if ctx.Response != nil {
			headers = ctx.Response.Header
		}
	default:
		if ctx.Headers != nil {
			headers = ctx.Headers
		}
	}

	if headers == nil {
		return errors.New("no headers available in context")
	}

	headers.Del(p.HeaderName)
	return nil
}

// Validate validates the policy configuration
func (p *RemoveHeaderPolicy) Validate() error {
	if p.HeaderName == "" {
		return errors.New("header name cannot be empty")
	}
	return nil
}

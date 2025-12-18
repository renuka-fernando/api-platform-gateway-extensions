package policies

import (
	"errors"
	"net/http"
)

// AddHeaderPolicy adds HTTP headers to request or response
type AddHeaderPolicy struct {
	// HeaderName is the name of the header to add
	HeaderName string
	// HeaderValue is the value of the header to add
	HeaderValue string
	// Flow specifies whether to apply on request or response
	Flow Flow
}

// NewAddHeaderPolicy creates a new AddHeaderPolicy
func NewAddHeaderPolicy(headerName, headerValue string, flow Flow) *AddHeaderPolicy {
	return &AddHeaderPolicy{
		HeaderName:  headerName,
		HeaderValue: headerValue,
		Flow:        flow,
	}
}

// Name returns the name of the policy
func (p *AddHeaderPolicy) Name() string {
	return "AddHeader"
}

// Execute applies the add header policy to the context
func (p *AddHeaderPolicy) Execute(ctx *PolicyContext) error {
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

	headers.Set(p.HeaderName, p.HeaderValue)
	return nil
}

// Validate validates the policy configuration
func (p *AddHeaderPolicy) Validate() error {
	if p.HeaderName == "" {
		return errors.New("header name cannot be empty")
	}
	if p.HeaderValue == "" {
		return errors.New("header value cannot be empty")
	}
	return nil
}

package policies

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// InterceptorPolicy calls an external service for custom processing
type InterceptorPolicy struct {
	// ServiceURL is the URL of the external interceptor service
	ServiceURL string
	// Timeout for the interceptor service call
	Timeout time.Duration
	// Flow specifies whether to apply on request or response
	Flow Flow
	// HTTPClient for making requests
	HTTPClient *http.Client
}

// InterceptorRequest represents the request sent to the interceptor service
type InterceptorRequest struct {
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body,omitempty"`
	Method  string              `json:"method,omitempty"`
	Path    string              `json:"path,omitempty"`
}

// InterceptorResponse represents the response from the interceptor service
type InterceptorResponse struct {
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body,omitempty"`
	Status  int                 `json:"status,omitempty"`
}

// NewInterceptorPolicy creates a new InterceptorPolicy
func NewInterceptorPolicy(serviceURL string, flow Flow) *InterceptorPolicy {
	return &InterceptorPolicy{
		ServiceURL: serviceURL,
		Timeout:    30 * time.Second,
		Flow:       flow,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Name returns the name of the policy
func (p *InterceptorPolicy) Name() string {
	return "Interceptor"
}

// Execute applies the interceptor policy to the context
func (p *InterceptorPolicy) Execute(ctx *PolicyContext) error {
	if ctx == nil {
		return errors.New("policy context cannot be nil")
	}

	var req *http.Request
	var headers http.Header
	var body string

	switch p.Flow {
	case RequestFlow:
		if ctx.Request == nil {
			return errors.New("request is nil")
		}
		req = ctx.Request
		headers = req.Header
		
		// Read body if present
		if req.Body != nil {
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				return fmt.Errorf("failed to read request body: %w", err)
			}
			body = string(bodyBytes)
			// Restore the body for subsequent handlers
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	case ResponseFlow:
		if ctx.Response == nil {
			return errors.New("response is nil")
		}
		headers = ctx.Response.Header
		
		// Read response body if present
		if ctx.Response.Body != nil {
			bodyBytes, err := io.ReadAll(ctx.Response.Body)
			if err != nil {
				return fmt.Errorf("failed to read response body: %w", err)
			}
			body = string(bodyBytes)
			// Restore the body for subsequent handlers
			ctx.Response.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	default:
		headers = ctx.Headers
	}

	if headers == nil {
		return errors.New("no headers available in context")
	}

	// Prepare interceptor request
	interceptorReq := InterceptorRequest{
		Headers: headers,
		Body:    body,
	}
	
	if req != nil {
		interceptorReq.Method = req.Method
		interceptorReq.Path = req.URL.Path
	}

	// Call interceptor service
	interceptorResp, err := p.callInterceptor(interceptorReq)
	if err != nil {
		return fmt.Errorf("interceptor service call failed: %w", err)
	}

	// Apply the response from interceptor
	for key, values := range interceptorResp.Headers {
		headers.Del(key)
		for _, value := range values {
			headers.Add(key, value)
		}
	}

	return nil
}

// callInterceptor makes an HTTP call to the interceptor service
func (p *InterceptorPolicy) callInterceptor(req InterceptorRequest) (*InterceptorResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", p.ServiceURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call interceptor service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("interceptor service returned status %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var interceptorResp InterceptorResponse
	if err := json.Unmarshal(respBody, &interceptorResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &interceptorResp, nil
}

// Validate validates the policy configuration
func (p *InterceptorPolicy) Validate() error {
	if p.ServiceURL == "" {
		return errors.New("service URL cannot be empty")
	}
	if p.Timeout <= 0 {
		return errors.New("timeout must be positive")
	}
	return nil
}

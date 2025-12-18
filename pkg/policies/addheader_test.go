package policies

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddHeaderPolicy_Name(t *testing.T) {
	policy := NewAddHeaderPolicy("X-Custom-Header", "test-value", RequestFlow)
	assert.Equal(t, "AddHeader", policy.Name())
}

func TestAddHeaderPolicy_Validate(t *testing.T) {
	tests := []struct {
		name        string
		headerName  string
		headerValue string
		wantErr     bool
	}{
		{
			name:        "valid policy",
			headerName:  "X-Custom-Header",
			headerValue: "test-value",
			wantErr:     false,
		},
		{
			name:        "empty header name",
			headerName:  "",
			headerValue: "test-value",
			wantErr:     true,
		},
		{
			name:        "empty header value",
			headerName:  "X-Custom-Header",
			headerValue: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := NewAddHeaderPolicy(tt.headerName, tt.headerValue, RequestFlow)
			err := policy.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddHeaderPolicy_Execute_Request(t *testing.T) {
	policy := NewAddHeaderPolicy("X-Custom-Header", "test-value", RequestFlow)
	
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	ctx := &PolicyContext{
		Request: req,
	}

	err := policy.Execute(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "test-value", req.Header.Get("X-Custom-Header"))
}

func TestAddHeaderPolicy_Execute_Response(t *testing.T) {
	policy := NewAddHeaderPolicy("X-Custom-Header", "test-value", ResponseFlow)
	
	resp := &http.Response{
		Header: make(http.Header),
	}
	ctx := &PolicyContext{
		Response: resp,
	}

	err := policy.Execute(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "test-value", resp.Header.Get("X-Custom-Header"))
}

func TestAddHeaderPolicy_Execute_NilContext(t *testing.T) {
	policy := NewAddHeaderPolicy("X-Custom-Header", "test-value", RequestFlow)
	err := policy.Execute(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context cannot be nil")
}

func TestAddHeaderPolicy_Execute_WithGenericHeaders(t *testing.T) {
	policy := NewAddHeaderPolicy("X-Custom-Header", "test-value", "")
	
	headers := make(http.Header)
	ctx := &PolicyContext{
		Headers: headers,
	}

	err := policy.Execute(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "test-value", headers.Get("X-Custom-Header"))
}

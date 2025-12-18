package policies

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveHeaderPolicy_Name(t *testing.T) {
	policy := NewRemoveHeaderPolicy("X-Custom-Header", RequestFlow)
	assert.Equal(t, "RemoveHeader", policy.Name())
}

func TestRemoveHeaderPolicy_Validate(t *testing.T) {
	tests := []struct {
		name       string
		headerName string
		wantErr    bool
	}{
		{
			name:       "valid policy",
			headerName: "X-Custom-Header",
			wantErr:    false,
		},
		{
			name:       "empty header name",
			headerName: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := NewRemoveHeaderPolicy(tt.headerName, RequestFlow)
			err := policy.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRemoveHeaderPolicy_Execute_Request(t *testing.T) {
	policy := NewRemoveHeaderPolicy("X-Remove-Me", RequestFlow)
	
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("X-Remove-Me", "should-be-removed")
	req.Header.Set("X-Keep-Me", "should-stay")
	
	ctx := &PolicyContext{
		Request: req,
	}

	err := policy.Execute(ctx)
	assert.NoError(t, err)
	assert.Empty(t, req.Header.Get("X-Remove-Me"))
	assert.Equal(t, "should-stay", req.Header.Get("X-Keep-Me"))
}

func TestRemoveHeaderPolicy_Execute_Response(t *testing.T) {
	policy := NewRemoveHeaderPolicy("X-Remove-Me", ResponseFlow)
	
	resp := &http.Response{
		Header: make(http.Header),
	}
	resp.Header.Set("X-Remove-Me", "should-be-removed")
	resp.Header.Set("X-Keep-Me", "should-stay")
	
	ctx := &PolicyContext{
		Response: resp,
	}

	err := policy.Execute(ctx)
	assert.NoError(t, err)
	assert.Empty(t, resp.Header.Get("X-Remove-Me"))
	assert.Equal(t, "should-stay", resp.Header.Get("X-Keep-Me"))
}

func TestRemoveHeaderPolicy_Execute_NilContext(t *testing.T) {
	policy := NewRemoveHeaderPolicy("X-Remove-Me", RequestFlow)
	err := policy.Execute(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context cannot be nil")
}

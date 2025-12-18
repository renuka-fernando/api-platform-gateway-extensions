package policies

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecutor_AddPolicy(t *testing.T) {
	executor := NewExecutor()
	
	policy := NewAddHeaderPolicy("X-Test", "value", RequestFlow)
	err := executor.AddPolicy(policy)
	assert.NoError(t, err)
	assert.Len(t, executor.GetPolicies(), 1)
}

func TestExecutor_AddPolicy_Nil(t *testing.T) {
	executor := NewExecutor()
	err := executor.AddPolicy(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "policy cannot be nil")
}

func TestExecutor_AddPolicy_InvalidPolicy(t *testing.T) {
	executor := NewExecutor()
	
	// Create an invalid policy (empty header name)
	policy := NewAddHeaderPolicy("", "value", RequestFlow)
	err := executor.AddPolicy(policy)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestExecutor_Execute(t *testing.T) {
	executor := NewExecutor()
	
	// Add multiple policies
	addPolicy := NewAddHeaderPolicy("X-Custom", "test-value", RequestFlow)
	removePolicy := NewRemoveHeaderPolicy("X-Remove", RequestFlow)
	
	err := executor.AddPolicy(addPolicy)
	assert.NoError(t, err)
	
	err = executor.AddPolicy(removePolicy)
	assert.NoError(t, err)
	
	// Execute policies
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("X-Remove", "should-be-removed")
	
	ctx := &PolicyContext{
		Request: req,
	}
	
	err = executor.Execute(ctx)
	assert.NoError(t, err)
	
	// Verify results
	assert.Equal(t, "test-value", req.Header.Get("X-Custom"))
	assert.Empty(t, req.Header.Get("X-Remove"))
}

func TestExecutor_Execute_NilContext(t *testing.T) {
	executor := NewExecutor()
	
	policy := NewAddHeaderPolicy("X-Test", "value", RequestFlow)
	err := executor.AddPolicy(policy)
	assert.NoError(t, err)
	
	err = executor.Execute(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context cannot be nil")
}

func TestExecutor_Clear(t *testing.T) {
	executor := NewExecutor()
	
	policy := NewAddHeaderPolicy("X-Test", "value", RequestFlow)
	err := executor.AddPolicy(policy)
	assert.NoError(t, err)
	assert.Len(t, executor.GetPolicies(), 1)
	
	executor.Clear()
	assert.Len(t, executor.GetPolicies(), 0)
}

func TestExecutor_ChainedExecution(t *testing.T) {
	executor := NewExecutor()
	
	// Chain multiple add header policies
	policy1 := NewAddHeaderPolicy("X-Header-1", "value-1", RequestFlow)
	policy2 := NewAddHeaderPolicy("X-Header-2", "value-2", RequestFlow)
	policy3 := NewAddHeaderPolicy("X-Header-3", "value-3", RequestFlow)
	
	executor.AddPolicy(policy1)
	executor.AddPolicy(policy2)
	executor.AddPolicy(policy3)
	
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	ctx := &PolicyContext{Request: req}
	
	err := executor.Execute(ctx)
	assert.NoError(t, err)
	
	assert.Equal(t, "value-1", req.Header.Get("X-Header-1"))
	assert.Equal(t, "value-2", req.Header.Get("X-Header-2"))
	assert.Equal(t, "value-3", req.Header.Get("X-Header-3"))
}

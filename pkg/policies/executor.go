package policies

import (
	"fmt"
)

// Executor manages the execution of multiple policies
type Executor struct {
	policies []Policy
}

// NewExecutor creates a new policy executor
func NewExecutor() *Executor {
	return &Executor{
		policies: make([]Policy, 0),
	}
}

// AddPolicy adds a policy to the executor
func (e *Executor) AddPolicy(policy Policy) error {
	if policy == nil {
		return fmt.Errorf("policy cannot be nil")
	}
	
	if err := policy.Validate(); err != nil {
		return fmt.Errorf("policy validation failed: %w", err)
	}
	
	e.policies = append(e.policies, policy)
	return nil
}

// Execute executes all policies in order
func (e *Executor) Execute(ctx *PolicyContext) error {
	if ctx == nil {
		return fmt.Errorf("policy context cannot be nil")
	}
	
	for _, policy := range e.policies {
		if err := policy.Execute(ctx); err != nil {
			return fmt.Errorf("policy %s execution failed: %w", policy.Name(), err)
		}
	}
	
	return nil
}

// GetPolicies returns the list of policies in the executor
func (e *Executor) GetPolicies() []Policy {
	return e.policies
}

// Clear removes all policies from the executor
func (e *Executor) Clear() {
	e.policies = make([]Policy, 0)
}

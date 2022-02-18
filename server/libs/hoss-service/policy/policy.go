package policy

import (
	"fmt"
)

// Policy is the top level policy object, consisting of a set of Statements
type Policy struct {
	Version    string
	Effect     string
	Statements []Statement
}

// Parse validates the Policy and returns the PolicyFilter that will apply this Policy to a message
func (policy *Policy) Parse() (PolicyFilter, error) {
	if policy.Version != "1" {
		return nil, fmt.Errorf("unsupported Policy Version value '%s'", policy.Version)
	}

	if policy.Effect == "" {
		policy.Effect = "OR"
	}

	if !contains(policy.Effect, []string{"AND", "OR"}) {
		return nil, fmt.Errorf("unsupported Effect value '%s'", policy.Effect)
	}

	stmts := make([]PolicyFilter, len(policy.Statements))
	for i := range policy.Statements {
		f, err := policy.Statements[i].Parse()
		if err != nil {
			return nil, err
		}

		stmts[i] = f
	}

	return LogicPolicies(policy.Effect, stmts)
}

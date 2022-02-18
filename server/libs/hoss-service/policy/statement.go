package policy

import "fmt"

// Statement is the mid-level policy object, consisting of a set of Conditions
type Statement struct {
	Id         string
	Effect     string
	Conditions []Condition
}

// Parse validates the Statement and returns the PolicyFilter that will apply this Statement to a message
func (stmt *Statement) Parse() (PolicyFilter, error) {
	if stmt.Effect == "" {
		stmt.Effect = "AND"
	}

	if !contains(stmt.Effect, []string{"AND", "OR"}) {
		return nil, fmt.Errorf("unsupported Effect value '%s'", stmt.Effect)
	}

	conds := make([]PolicyFilter, len(stmt.Conditions))
	for i := range stmt.Conditions {
		f, err := stmt.Conditions[i].Parse()
		if err != nil {
			return nil, err
		}

		conds[i] = f
	}

	return LogicPolicies(stmt.Effect, conds)
}

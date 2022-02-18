package policy

import (
	"fmt"

	"github.com/ghodss/yaml"
)

// MessageInformation is all of the information about a notification message that is required for a PolicyFilter to make a policy decision about the message
type MessageInformation struct {
	EventOperation string

	ObjectKey  string
	ObjectSize int

	ObjectMetadata map[string]string
}

// PolicyFilter is the signature of the a filter function that applied a parsed policy to a message
// Return is true if the message meets the policy requirements
type PolicyFilter func(msg *MessageInformation) (bool, error)

// DefaultOpenPolicy is the policy document that allows all messages
const DefaultOpenPolicy = "{\"Version\":\"1\",\"Statements\":[]}"

// DefaultOpenPolicyFilter is the policy filter that allows all messages
func DefaultOpenPolicyFilter(msg *MessageInformation) (bool, error) {
	return true, nil
}

// Parse takes a policy document, parses it, and returns the PolicyFilter that will make decisions about messages
func Parse(policyDoc string) (PolicyFilter, error) {
	policyObj := &Policy{}
	if err := yaml.Unmarshal([]byte(policyDoc), &policyObj); err != nil {
		return nil, fmt.Errorf("Could not parse policy document: %w", err)
	}

	return policyObj.Parse()
}

// LogicPolicies is a common method for handling ORing or ANDing a list of policies together
func LogicPolicies(effect string, policies []PolicyFilter) (PolicyFilter, error) {
	// shortcut is used to shortcut the boolean operations
	// if shortcut == true then the policies are ORed together
	// if shortcut == false then the policies are ANDed together
	var shortcut bool
	switch effect {
	case "OR":
		shortcut = true
	case "AND":
		shortcut = false
	default:
		return nil, fmt.Errorf("unsupported Effect value '%s'", effect)
	}

	return func(msg *MessageInformation) (bool, error) {
		if len(policies) == 0 {
			return true, nil
		}

		for _, policy := range policies {
			val, err := policy(msg)
			if err != nil {
				return false, err
			}

			if val == shortcut {
				return shortcut, nil
			}
		}

		return !shortcut, nil
	}, nil
}

// contains returns true if the needle is found in the haystack
func contains(needle string, haystack []string) bool {
	for _, hay := range haystack {
		if hay == needle {
			return true
		}
	}
	return false
}

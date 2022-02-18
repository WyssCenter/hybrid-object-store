package policy

import (
	"fmt"
	"testing"
)

type TestMessageInformation struct {
	Message  *MessageInformation
	Expected bool
}

func testPolicy(t *testing.T, policyDoc string, msg *MessageInformation) (bool, error) {
	policyFilter, err := Parse(policyDoc)
	if err != nil {
		return false, fmt.Errorf("parse failed: %w", err)
	}

	passed, err := policyFilter(msg)
	if err != nil {
		return false, fmt.Errorf("policy filter failed: %w", err)
	}

	return passed, nil
}

func TestDefaultPolicy(t *testing.T) {
	policy := "{\"Version\":\"1\",\"Statements\":[]}"
	msg := &MessageInformation{
		EventOperation: "PUT",
		ObjectKey:      "file.ext",
		ObjectSize:     1024, // bytes
		ObjectMetadata: map[string]string{},
	}

	passed, err := testPolicy(t, policy, msg)
	if err == nil && !passed {
		err = fmt.Errorf("policy filter didn't evaluate to true")
	}
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestSimplePolicyPass(t *testing.T) {
	policy := `{
		"Version": "1",
		"Statements": [{
			"Id": "CheckFileKey",
			"Conditions": [{
				"Left": "object:key",
				"Right": "file.key",
				"Operator": "==",
			}],
		}],
	}`
	msg := &MessageInformation{
		EventOperation: "PUT",
		ObjectKey:      "file.key",
		ObjectSize:     1024, // bytes
		ObjectMetadata: map[string]string{},
	}

	passed, err := testPolicy(t, policy, msg)
	if err == nil && !passed {
		err = fmt.Errorf("policy filter didn't evaluate to true")
	}
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestMultipleStatements(t *testing.T) {
	policy := `{
		"Version": "1",
		"Statements": [{
			"Id": "CheckRawFiles",
			"Conditions": [{
				"Left": "object:key",
				"Right": "*.raw",
				"Operator": "!=",
			}],
		},{
			"Id": "CheckLargeFiles",
			"Conditions": [{
				"Left": "object:size",
				"Right": 4096,
				"Operator": "<",
			}],
		}],
	}`

	tmis := []TestMessageInformation{
		{
			// Pass both
			Message: &MessageInformation{
				EventOperation: "PUT",
				ObjectKey:      "file.key",
				ObjectSize:     1024, // bytes
				ObjectMetadata: map[string]string{},
			},
			Expected: true,
		},
		{
			// Fail CheckLargeFiles
			Message: &MessageInformation{
				EventOperation: "PUT",
				ObjectKey:      "file.key",
				ObjectSize:     8192, // bytes
				ObjectMetadata: map[string]string{},
			},
			Expected: true,
		},
		{
			// fail CheckRawFiles
			Message: &MessageInformation{
				EventOperation: "PUT",
				ObjectKey:      "file.raw",
				ObjectSize:     1024, // bytes
				ObjectMetadata: map[string]string{},
			},
			Expected: true,
		},
		{
			// Fail both
			Message: &MessageInformation{
				EventOperation: "PUT",
				ObjectKey:      "file.raw",
				ObjectSize:     8192, // bytes
				ObjectMetadata: map[string]string{},
			},
			Expected: false,
		},
	}

	for i, tmi := range tmis {
		passed, err := testPolicy(t, policy, tmi.Message)
		if err == nil && passed != tmi.Expected {
			err = fmt.Errorf("message %d policy filter didn't evaluate to %t", i+1, tmi.Expected)
		}
		if err != nil {
			t.Fatalf(err.Error())
		}
	}
}

func TestMultipleCondition(t *testing.T) {
	policy := `{
		"Version": "1",
		"Statements": [{
			"Id": "CheckLargeRawFiles",
			"Conditions": [{
				"Left": "object:key",
				"Right": "*.raw",
				"Operator": "!=",
			},{
				"Left": "object:size",
				"Right": 4096,
				"Operator": "<",
			}],
		}],
	}`

	tmis := []TestMessageInformation{
		{
			// Pass both
			Message: &MessageInformation{
				EventOperation: "PUT",
				ObjectKey:      "file.key",
				ObjectSize:     1024, // bytes
				ObjectMetadata: map[string]string{},
			},
			Expected: true,
		},
		{
			// Fail size check
			Message: &MessageInformation{
				EventOperation: "PUT",
				ObjectKey:      "file.key",
				ObjectSize:     8192, // bytes
				ObjectMetadata: map[string]string{},
			},
			Expected: false,
		},
		{
			// fail raw check
			Message: &MessageInformation{
				EventOperation: "PUT",
				ObjectKey:      "file.raw",
				ObjectSize:     1024, // bytes
				ObjectMetadata: map[string]string{},
			},
			Expected: false,
		},
		{
			// Fail both
			Message: &MessageInformation{
				EventOperation: "PUT",
				ObjectKey:      "file.raw",
				ObjectSize:     8192, // bytes
				ObjectMetadata: map[string]string{},
			},
			Expected: false,
		},
	}

	for i, tmi := range tmis {
		passed, err := testPolicy(t, policy, tmi.Message)
		if err == nil && passed != tmi.Expected {
			err = fmt.Errorf("message %d policy filter didn't evaluate to %t", i+1, tmi.Expected)
		}
		if err != nil {
			t.Fatalf(err.Error())
		}
	}
}

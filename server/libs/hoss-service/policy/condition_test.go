package policy

import (
	"testing"
)

func testCond(t *testing.T, cond *Condition, msg *MessageInformation) {
	policyFilter, err := cond.Parse()
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	result, err := policyFilter(msg)
	if err != nil {
		t.Fatalf("policy filter failed: %v", err)
	}

	if result == false {
		t.Fatalf("policy filter didn't evaluate to true")
	}
}

func TestEventOperation(t *testing.T) {
	msg := &MessageInformation{
		EventOperation: "PUT",
		ObjectKey:      "file.ext",
		ObjectSize:     1024, // bytes
		ObjectMetadata: map[string]string{},
	}

	cond := &Condition{
		Left:     "event:operation",
		Right:    []byte(`"PUT"`),
		Operator: "==",
	}

	testCond(t, cond, msg)
}

func TestObjectKey(t *testing.T) {
	msg := &MessageInformation{
		EventOperation: "PUT",
		ObjectKey:      "file.ext",
		ObjectSize:     1024, // bytes
		ObjectMetadata: map[string]string{},
	}

	cond := &Condition{
		Left:     "object:key",
		Right:    []byte(`"*.ext"`),
		Operator: "==",
	}

	testCond(t, cond, msg)
}

func TestObjectSize(t *testing.T) {
	msg := &MessageInformation{
		EventOperation: "PUT",
		ObjectKey:      "file.ext",
		ObjectSize:     1024, // bytes
		ObjectMetadata: map[string]string{},
	}

	cond := &Condition{
		Left:     "object:size",
		Right:    []byte(`2048`),
		Operator: "<",
	}

	testCond(t, cond, msg)
}

func TestObjectMetadataKeyExists(t *testing.T) {
	msg := &MessageInformation{
		EventOperation: "PUT",
		ObjectKey:      "file.ext",
		ObjectSize:     1024, // bytes
		ObjectMetadata: map[string]string{
			"test": "test",
		},
	}

	cond := &Condition{
		Left:     "object:metadata",
		Right:    []byte(`"test"`),
		Operator: "has",
	}

	testCond(t, cond, msg)
}

func TestObjectMetadataValueSet(t *testing.T) {
	msg := &MessageInformation{
		EventOperation: "PUT",
		ObjectKey:      "file.ext",
		ObjectSize:     1024, // bytes
		ObjectMetadata: map[string]string{
			"test": "",
		},
	}

	cond := &Condition{
		Left:     "object:metadata:test",
		Right:    []byte(`""`),
		Operator: "==",
	}

	testCond(t, cond, msg)
}

func TestObjectMetadataValueGlob(t *testing.T) {
	msg := &MessageInformation{
		EventOperation: "PUT",
		ObjectKey:      "file.ext",
		ObjectSize:     1024, // bytes
		ObjectMetadata: map[string]string{
			"test": "test",
		},
	}

	cond := &Condition{
		Left:     "object:metadata:test",
		Right:    []byte(`"foo*"`),
		Operator: "!=",
	}

	testCond(t, cond, msg)
}

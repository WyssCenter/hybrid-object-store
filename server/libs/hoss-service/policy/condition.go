package policy

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gobwas/glob"
)

// Condition is a single policy condition that applies an operator to two operands
type Condition struct {
	// Left operand
	Left string

	// Right operand - type can be either number or string, so we defer parsing until later
	Right json.RawMessage

	// Fields to hold parsed value from Right operand
	rightNum float64
	rightStr string

	// Operator that will be applied to the two operands
	Operator string
}

// validateLeft ensures that the Left operand is a valid value
func (cond *Condition) validateLeft() error {
	if cond.Left == "" {
		return fmt.Errorf("left value is required: %s", cond.Left)
	}
	switch cond.Left {
	case "event:operation":
		return nil
	case "object:key":
		return nil
	case "object:size":
		return nil
	case "object:metadata":
		return nil
	default:
		if strings.HasPrefix(cond.Left, "object:metadata:") {
			return nil
		} else {
			return fmt.Errorf("invalid Left value: %s", cond.Left)
		}
	}
}

// lookupLeft will return the value in the MessageInformation that is the target of the Left operand
func (cond *Condition) lookupLeft(msg *MessageInformation) (interface{}, error) {
	switch cond.Left {
	case "event:operation":
		switch msg.EventOperation {
		case "s3:ObjectCreated:Put",
			"s3:ObjectCreated:Copy",
			"s3:ObjectCreated:CompleteMultipartUpload",
			"ObjectCreated:Put", // AWS doesn't include the 's3:' prefix
			"ObjectCreated:Copy",
			"ObjectCreated:CompleteMultipartUpload":
			return "PUT", nil
		case "s3:ObjectRemoved:Delete",
			"ObjectRemoved:Delete",
			"ObjectRemoved:DeleteMarkerCreated",
			"s3:ObjectRemoved:DeleteMarkerCreated":
			return "DELETE", nil
		default:
			return "UNSUPPORTED", nil
		}

	case "object:key":
		return msg.ObjectKey, nil
	case "object:size":
		return float64(msg.ObjectSize), nil
	case "object:metadata":
		return msg.ObjectMetadata, nil
	default:
		if strings.HasPrefix(cond.Left, "object:metadata:") {
			key := cond.Left[len("object:metadata:"):len(cond.Left)]
			value, ok := msg.ObjectMetadata[key]
			if !ok {
				return nil, fmt.Errorf("metadata key %s doesn't exist", key)
			} else {
				return value, nil
			}
		} else {
			return nil, fmt.Errorf("invalid Left value: %s", cond.Left)
		}
	}
}

// validateRight ensures that the Right operand is a valid value
// This is where the Right operand is parsed from a RawMessage, based on the type of the Left operand
func (cond *Condition) validateRight() error {
	if len(cond.Right) == 0 {
		return fmt.Errorf("Right operand is required")
	}

	if cond.Left == "event:operation" ||
		cond.Left == "object:key" ||
		cond.Left == "object:metadata" ||
		strings.HasPrefix(cond.Left, "object:metadata:") {
		// string Right operand

		str := string(cond.Right)
		if str[0] != '"' || str[len(str)-1] != '"' {
			return fmt.Errorf("Left operand expects a string Right operand")
		}

		str = str[1 : len(str)-1] // remove the quote characters

		_, err := glob.Compile(str)
		if err != nil {
			return fmt.Errorf("problem compiling glob (%q): %v", str, err)
		}

		cond.rightStr = str
	} else if cond.Left == "object:size" {
		// num Right operand

		str := string(cond.Right)
		num, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return fmt.Errorf("Left operand expects a numeric Right operand")
		}

		cond.rightStr = str
		cond.rightNum = num
	} else {
		return fmt.Errorf("unknown type for Left %s", cond.Left)
	}

	return nil
}

// validateOperator ensures that the Operator is a valid operator for the given operands
func (cond *Condition) validateOperator() error {
	if cond.Left == "event:operation" ||
		cond.Left == "object:key" ||
		strings.HasPrefix(cond.Left, "object:metadata:") {
		// string value
		if !contains(cond.Operator, []string{"==", "!="}) {
			return fmt.Errorf("cannot apply operator %q to string values", cond.Operator)
		}
	} else if cond.Left == "object:size" {
		// num value
		if !contains(cond.Operator, []string{"==", "!=", "<", "<=", ">", ">="}) {
			return fmt.Errorf("cannot apply operator %q to num values", cond.Operator)
		}
	} else if cond.Left == "object:metadata" {
		// map value
		if !contains(cond.Operator, []string{"has"}) {
			return fmt.Errorf("cannot apply operator %q to map key check", cond.Operator)
		}
	} else {
		return fmt.Errorf("unknown type for Left %s", cond.Left)
	}

	return nil
}

// applyOperator applied the Operator to the Operands (both the looked up left value and the previously parsed right value)
func (cond *Condition) applyOperator(left interface{}) bool {
	// Note: this function assumes that the validate* functions have been called, as
	//       error checking / handling is not done in this function
	var passes bool

	switch left.(type) {
	case string:
		g, _ := glob.Compile(cond.rightStr)
		passes = g.Match(left.(string))

		if cond.Operator == "!=" {
			passes = !passes
		}
	case float64:
		leftValue := left.(float64)
		rightValue := cond.rightNum

		switch cond.Operator {
		case "<":
			passes = leftValue < rightValue
		case ">":
			passes = leftValue > rightValue
		case "<=":
			passes = leftValue <= rightValue
		case ">=":
			passes = leftValue >= rightValue
		case "==":
			passes = leftValue == rightValue
		case "!=":
			passes = leftValue != rightValue
		}
	case map[string]string:
		leftValue := left.(map[string]string)
		_, passes = leftValue[cond.rightStr]
	}

	return passes
}

// Parse validates the Condition and returns the PolicyFilter that will apply this condition to a message
func (cond *Condition) Parse() (PolicyFilter, error) {
	if err := cond.validateLeft(); err != nil {
		return nil, err
	}

	if err := cond.validateRight(); err != nil {
		return nil, err
	}

	if err := cond.validateOperator(); err != nil {
		return nil, err
	}

	return func(msg *MessageInformation) (bool, error) {
		leftValue, err := cond.lookupLeft(msg)
		if err != nil {
			return false, err
		}

		return cond.applyOperator(leftValue), nil

	}, nil
}

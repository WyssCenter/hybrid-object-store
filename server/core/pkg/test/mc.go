package test

import (
	"os/exec"
	"testing"
)

// MinioPolicyExists returns true if a policy exists
func MinioPolicyExists(t *testing.T, policyName string) bool {
	cmd := exec.Command("mc", "admin", "policy", "info", "hoss-default", policyName, "--json")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

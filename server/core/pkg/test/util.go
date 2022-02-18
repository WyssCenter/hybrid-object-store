package test

import (
	"fmt"
	"testing"

	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// AssertEqual is a helper function to check 2 values
func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

// DefaultBucketDir returns the default bucket dir
func DefaultBucketDir(t *testing.T, defaultBucket string) string {
	nasDir, err := homedir.Expand("~/.hoss/data/nas")
	if err != nil {
		t.Fatalf("Failed to get directory for default bucket")
	}

	p, err := filepath.Abs(fmt.Sprintf("%s/%s", nasDir, defaultBucket))
	if err != nil {
		t.Fatalf("Failed to get directory for default bucket")
	}
	return p
}

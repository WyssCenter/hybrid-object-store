package api

import (
	"os"

	"github.com/pkg/errors"

	yaml "github.com/ghodss/yaml"
)

var (
	// ErrAccountNotAllowed is raised when an account is not on the allowlist
	ErrAccountNotAllowed = errors.New("User's email is not allowlisted")
)

// AccountInformation contains the information needed to verify a user's account
type AccountInformation struct {
	Email string
}

// AccountAllowlist defines the allowlist YAML input data
type AccountAllowlist struct {
	EnforceAllowlist bool     `yaml:"enforce_allowlist,omitempty"`
	Emails           []string `yaml:"emails,omitempty"`
}

// LoadAllowlist loads the allowlist information from the ACCOUNT_ALLOWLIST environment variable
func LoadAllowlist() (*AccountAllowlist, error) {
	allowlist := &AccountAllowlist{}

	data := os.Getenv("ACCOUNT_ALLOWLIST")
	if data != "" {
		err := yaml.Unmarshal([]byte(data), &allowlist)
		if err != nil {
			return nil, errors.Wrap(err, "Cannot parse allowlist data")
		}
	}

	return allowlist, nil
}

// CheckAllowlist checks to see if the given account has been allowlisted
func (allowlist *AccountAllowlist) CheckAllowlist(account AccountInformation) error {
	if !allowlist.EnforceAllowlist {
		return nil
	}

	for _, v := range allowlist.Emails {
		if account.Email == v {
			return nil
		}
	}

	return ErrAccountNotAllowed
}

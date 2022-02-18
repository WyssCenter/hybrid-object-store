package test

import (
	"bufio"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

//LoadEnvFile Loads the env file containing the password for minio
func LoadEnvFile(path string) error {

	// Use homedir package so this works easily on mac/linux (useful when developing)
	absPath, err := homedir.Expand(path)
	if err != nil {
		return err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var envvars []string
	for scanner.Scan() {
		envvars = append(envvars, scanner.Text())
	}
	file.Close()

	// and then a loop iterates through
	// and prints each of the slice values.
	for _, ev := range envvars {
		evSplit := strings.Split(ev, "=")
		err = os.Setenv(evSplit[0], evSplit[1])
		if err != nil {
			return err
		}
	}

	return nil
}

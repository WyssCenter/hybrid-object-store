package store

import (
	"bytes"
	"fmt"
	"strings"

	"text/template"

	"github.com/gigantum/hoss-core/pkg/database"
	"github.com/pkg/errors"
)

type statementData struct {
	Bucket      string
	DatasetName string
}

// RenderTemplate is a utility function that combines templates and permissions into a full polcy
func RenderTemplate(policyTemplate string, readTemplate string,
	readWriteTemplate string, denyTemplate string, permissions []*database.Permission) (string, error) {

	var statements = []string{}

	if len(permissions) > 0 {
		// You have valid permissions that should be rendered
		readTmpl := template.New("r")
		readTmpl, err := readTmpl.Parse(readTemplate)
		if err != nil {
			return "", errors.Wrap(err, "Failed to parse read statement template")
		}

		readWriteTmpl := template.New("rw")
		readWriteTmpl, err = readWriteTmpl.Parse(readWriteTemplate)
		if err != nil {
			return "", errors.Wrap(err, "Failed to parse readwrite statement template")
		}

		// append read-write statements after read-only statements to ensure that read-only statements
		// can't overwrite read-write statements for the same resource
		var readStatements = []string{}
		var readWriteStatements = []string{}
		for _, perm := range permissions {
			var tmplStr bytes.Buffer
			s := &statementData{
				Bucket:      perm.Dataset.Namespace.BucketName,
				DatasetName: perm.Dataset.Name,
			}
			switch access := perm.Permission; access {
			case database.PERM_READ:
				// Render template to string and save in list
				err = readTmpl.Execute(&tmplStr, s)
				if err != nil {
					return "", errors.Wrap(err, "Failed to render policy")
				}
				readStatements = append(readStatements, tmplStr.String())

			case database.PERM_READ_WRITE:
				// Render template to string and save in list
				err = readWriteTmpl.Execute(&tmplStr, s)
				if err != nil {
					return "", errors.Wrap(err, "Failed to render policy")
				}
				readWriteStatements = append(readWriteStatements, tmplStr.String())

			default:
				return "", errors.New(fmt.Sprintf("Unsupported permission type: %v", &perm.Permission))

			}
		}
		statements = append(readStatements, readWriteStatements...)
	} else {
		// You have no permissions so just the deny policy should be used
		denyTmpl := template.New("deny")
		denyTmpl, err := denyTmpl.Parse(denyTemplate)
		if err != nil {
			return "", errors.Wrap(err, "Failed to parse readwrite statement template")
		}

		var tmplStr bytes.Buffer
		s := &statementData{}

		err = denyTmpl.Execute(&tmplStr, s)
		if err != nil {
			return "", errors.Wrap(err, "Failed to render policy")
		}
		statements = append(statements, tmplStr.String())
	}

	statementsStr := strings.Join(statements[:], ",\n")
	output := strings.ReplaceAll(policyTemplate, "{{STATEMENTS}}", statementsStr)

	return output, nil
}

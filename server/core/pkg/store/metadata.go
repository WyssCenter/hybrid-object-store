package store

import (
	"fmt"
	"time"
)

// MetadataFile is a struct representing a dataset metadata file
type MetadataFile struct {
	DatasetName string `yaml:"dataset_name"`
	CreatedOn   string `yaml:"created_on"`
}

// Key returns the filename for the metadata file within the root bucket
func (m *MetadataFile) Key() string {
	return fmt.Sprintf("%s/.dataset.yaml", m.DatasetName)
}

// NewMetadataFile creates a new metadata file struct
func NewMetadataFile(datasetName string) *MetadataFile {
	return &MetadataFile{
		DatasetName: datasetName,
		CreatedOn:   time.Now().Format(time.RFC3339),
	}
}

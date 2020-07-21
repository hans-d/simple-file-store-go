package filestore

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// Marshaler is used to read/write the data
type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	GetFileExtension() string
}

// JSONMarshaler uses the JSON file format.
type JSONMarshaler struct {
}

func (m JSONMarshaler) Marshal(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "\t")
}

func (m JSONMarshaler) Unmarshal(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}

func (m JSONMarshaler) GetFileExtension() string {
	return ".json"
}

// JSONMarshaler uses the YAML file format.
type YAMLMarshaler struct {
}

func (m YAMLMarshaler) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (m YAMLMarshaler) Unmarshal(b []byte, v interface{}) error {
	return yaml.Unmarshal(b, v)
}

func (m YAMLMarshaler) GetFileExtension() string {
	return ".yml"
}

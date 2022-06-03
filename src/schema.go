package src

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Label struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

type Repository struct {
	Owner  string  `yaml:"owner"`
	Name   string  `yaml:"name"`
	Labels []Label `yaml:"labels"`
}

type Schema struct {
	Repositories []Repository `yaml:"repositories"`
}

func LoadSchema(path string) (*Schema, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	schema := &Schema{}
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(schema)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

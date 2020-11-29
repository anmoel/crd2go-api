package convert

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	// "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

// CustomResourceDefinition describes the structure of a CRD file
type CustomResourceDefinition struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

// Metadata conatins the important fields of the CustomResourceDefinition metadata block
type Metadata struct {
	Name string `yaml:"name"`
}

// Spec conatins the important fields of the CustomResourceDefinition specs block
type Spec struct {
	Group      string     `yaml:"group"`
	Names      Names      `yaml:"names"`
	Scope      string     `yaml:"scope"`
	Validation Validation `yaml:"validation"`
	Version    string     `yaml:"version"`
}

// Names contains the crd name informations
type Names struct {
	Categories []string `yaml:"categories"`
	Kind       string   `yaml:"kind"`
	Plural     string   `yaml:"plural"`
	ShortNames []string `yaml:"shortNames"`
	Singular   string   `yaml:"singular"`
}

// Validation conatins the validation block of the CustomResourceDefinition
type Validation struct {
	OpenAPIV3Schema OpenAPIV3Schema `yaml:"openAPIV3Schema"`
}

// OpenAPIV3Schema conatins the important fields of the openAPIV3Schema block from the CustomResourceDefinition
type OpenAPIV3Schema struct {
	Description          string                     `yaml:"description,omitempty"`
	Properties           map[string]OpenAPIV3Schema `yaml:"properties,omitempty"`
	Required             []string                   `yaml:"required,omitempty"`
	Type                 string                     `yaml:"type"`
	Pattern              string                     `yaml:"pattern,omitempty"`
	Items                *OpenAPIV3Schema           `yaml:"items,omitempty"`
	OneOf                []OpenAPIV3Schema          `yaml:"oneOf,omitempty"`
	AnyOf                []OpenAPIV3Schema          `yaml:"anyOf,omitempty"`
	Not                  *OpenAPIV3Schema           `yaml:"not,omitempty"`
	AdditionalProperties *OpenAPIV3Schema           `yaml:"additionalProperties,omitempty"`
}

func readCrdFile(filepath string) (*CustomResourceDefinition, error) {
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var crd CustomResourceDefinition
	if err = yaml.Unmarshal(yamlFile, &crd); err != nil {
		return nil, err
	}

	return &crd, nil
}

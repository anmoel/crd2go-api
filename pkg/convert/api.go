package convert

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/anmoel/crd2go-api/pkg/templates"
)

type templateOptionsGroupversionInfo struct {
	Group   string
	Version string
}

type blockProperty struct {
	Name string
	Type string
	JSON string
}

type templateOptionsSpecStatusBlock struct {
	Kind        string
	Properties  []blockProperty
	Description string
}

var blockProperties map[string]*OpenAPIV3Schema

func createGroupversionInfoGoFile(filePath string, templateOptions *templateOptionsGroupversionInfo) error {
	return createTemplateFile(filePath, templates.TemplateGroupversionInfoGo, templateOptions)
}

func createTypesGoFile(filePath string, crd *CustomResourceDefinition) error {
	folderPath := filepath.Dir(filePath)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	fileTemplate, err := template.New("go-tmpl").Parse(templates.TemplateTypesGo)
	if err != nil {
		return err
	}
	err = fileTemplate.Execute(file, crd.Spec)
	if err != nil {
		return err
	}

	blockProperties = make(map[string]*OpenAPIV3Schema)
	//spec
	var specProperties []blockProperty
	specTemplate, err := template.New("spec-tmpl").Parse(templates.TemplateSpecBlock)
	if err != nil {
		return err
	}
	for key, value := range crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties {
		var jsonDefinition string
		if contains(crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Required, key) {
			jsonDefinition = fmt.Sprintf("`json:\"%s\"`", key)
		} else {
			jsonDefinition = fmt.Sprintf("`json:\"%s,omitempty\"`", key)
		}
		specProperties = append(specProperties, blockProperty{
			Name: strings.Title(key),
			Type: getPropertyTypeAndCreateNewBlock(key, value, folderPath),
			JSON: jsonDefinition,
		})
	}
	err = specTemplate.Execute(file, &templateOptionsSpecStatusBlock{
		Kind:        crd.Spec.Names.Kind,
		Properties:  specProperties,
		Description: crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Description,
	})
	if err != nil {
		return err
	}

	//status
	var statusProperties []blockProperty
	statusTemplate, err := template.New("status-tmpl").Parse(templates.TemplateStatusBlock)
	if err != nil {
		return err
	}
	for key, value := range crd.Spec.Validation.OpenAPIV3Schema.Properties["status"].Properties {
		var jsonDefinition string
		if contains(crd.Spec.Validation.OpenAPIV3Schema.Properties["status"].Required, key) {
			jsonDefinition = fmt.Sprintf("`json:\"%s\"`", key)
		} else {
			jsonDefinition = fmt.Sprintf("`json:\"%s,omitempty\"`", key)
		}
		statusProperties = append(statusProperties, blockProperty{
			Name: strings.Title(key),
			Type: getPropertyTypeAndCreateNewBlock(key, value, folderPath),
			JSON: jsonDefinition,
		})
	}
	if err := statusTemplate.Execute(file, &templateOptionsSpecStatusBlock{
		Kind:        crd.Spec.Names.Kind,
		Properties:  statusProperties,
		Description: crd.Spec.Validation.OpenAPIV3Schema.Properties["status"].Description,
	}); err != nil {
		return err
	}
	//additional blocks
	for {
		for propertyName, propertyValue := range blockProperties {
			var chieldProperties []blockProperty
			blockTemplate, err := template.New("status-tmpl").Parse(templates.TemplateBlock)
			if err != nil {
				return err
			}
			for key, value := range propertyValue.Properties {
				var jsonDefinition string
				if contains(propertyValue.Required, key) {
					jsonDefinition = fmt.Sprintf("`json:\"%s\"`", key)
				} else {
					jsonDefinition = fmt.Sprintf("`json:\"%s,omitempty\"`", key)
				}
				chieldProperties = append(chieldProperties, blockProperty{
					Name: strings.Title(key),
					Type: getPropertyTypeAndCreateNewBlock(key, value, folderPath),
					JSON: jsonDefinition,
				})
			}
			if err := blockTemplate.Execute(file, &templateOptionsSpecStatusBlock{
				Kind:        propertyName,
				Properties:  chieldProperties,
				Description: propertyValue.Description,
			}); err != nil {
				return err
			}
			delete(blockProperties, propertyName)
		}
		if len(blockProperties) == 0 {
			break
		}
	}

	file.Close()
	return nil
}

func createTemplateFile(filePath string, templateString string, templateOptions interface{}) error {
	template, err := template.New("go-tmpl").Parse(templateString)
	if err != nil {
		return err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	err = template.Execute(file, templateOptions)
	if err != nil {
		return err
	}
	file.Close()

	return nil
}

func getPropertyTypeAndCreateNewBlock(key string, value OpenAPIV3Schema, folderPath string) string {
	switch value.Type {
	case "array":
		kind := trimSuffix(strings.Title(key), "s")
		if value.Items.Type == "object" {
			exists, _ := propertyBlockExists(folderPath, kind)
			if !exists {
				blockProperties[kind] = value.Items
			}
			return fmt.Sprintf("[]%s", kind)
		}
		return fmt.Sprintf("[]%s", kind)

	case "object":
		kind := strings.Title(key)
		if value.AdditionalProperties != nil {
			return fmt.Sprintf("map[string]%s", value.AdditionalProperties.Type)
		}
		exists, _ := propertyBlockExists(folderPath, kind)
		if !exists {
			blockProperties[strings.Title(key)] = &value
		}
		return kind
	case "boolean":
		return "bool"
	case "integer":
		return "int"
	default:
		return value.Type
	}
}

func propertyBlockExists(folderPath string, kind string) (bool, error) {
	blockExists := false
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return blockExists, err
	}
	for _, file := range files {
		read, err := ioutil.ReadFile(file.Name())
		if err != nil {
			return blockExists, err
		}
		if strings.Contains(string(read), fmt.Sprintf("type %s struct {", kind)) {
			blockExists = true
		}
	}
	return false, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

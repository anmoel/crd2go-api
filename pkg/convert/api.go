package convert

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/anmoel/crd2go-api/pkg/templates"
)

type templateOptionsGroupversionInfo struct {
	Group   string
	Version string
	License string
}
type templateOptionsTypes struct {
	CRD     Spec
	License string
	Spec    bool
	Status  bool
}

type blockProperty struct {
	Name string
	Type string
	JSON string
}

type templateOptionsSpecStatusBlock struct {
	Kind       string
	Properties []blockProperty
}

type templateOptionsBlock struct {
	Kind        string
	Properties  []blockProperty
	Description string
}

var blockProperties map[string]*OpenAPIV3Schema

func createGroupversionInfoGoFile(filePath string, templateOptions *templateOptionsGroupversionInfo) error {
	return createTemplateFile(filePath, templates.TemplateGroupversionInfoGo, templateOptions)
}

func createTypesGoFile(filePath string, crd *CustomResourceDefinition, license string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	fileTemplate, err := template.New("go-tmpl").Parse(templates.TemplateTypesGo)
	if err != nil {
		return err
	}
	err = fileTemplate.Execute(file, templateOptionsTypes{
		License: license,
		CRD:     crd.Spec,
		Spec:    crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties != nil,
		Status:  crd.Spec.Validation.OpenAPIV3Schema.Properties["status"].Properties != nil,
	})
	if err != nil {
		return err
	}

	blockProperties = make(map[string]*OpenAPIV3Schema)
	//spec
	if crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties != nil {
		var specProperties []blockProperty
		specTemplate, err := template.New("spec-tmpl").Parse(templates.TemplateSpecBlock)
		if err != nil {
			return err
		}
		for key, value := range crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties {
			var jsonDefinition string
			required := true
			if contains(crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Required, key) {
				jsonDefinition = fmt.Sprintf("`json:\"%s\"`", key)
			} else {
				jsonDefinition = fmt.Sprintf("`json:\"%s,omitempty\"`", key)
				required = false
			}
			propertyType, err := getPropertyTypeAndCreateNewBlock(key, value, filePath, required)
			if err != nil {
				return err
			}
			specProperties = append(specProperties, blockProperty{
				Name: strings.Title(key),
				Type: propertyType,
				JSON: jsonDefinition,
			})
		}
		sort.Slice(specProperties, func(i, j int) bool { return specProperties[i].Name < specProperties[j].Name })
		err = specTemplate.Execute(file, &templateOptionsSpecStatusBlock{
			Kind:       crd.Spec.Names.Kind,
			Properties: specProperties,
		})
		if err != nil {
			return err
		}
	}

	//status
	if crd.Spec.Validation.OpenAPIV3Schema.Properties["status"].Properties != nil {
		var statusProperties []blockProperty
		statusTemplate, err := template.New("status-tmpl").Parse(templates.TemplateStatusBlock)
		if err != nil {
			return err
		}
		for key, value := range crd.Spec.Validation.OpenAPIV3Schema.Properties["status"].Properties {
			var jsonDefinition string
			required := true
			if contains(crd.Spec.Validation.OpenAPIV3Schema.Properties["status"].Required, key) {
				jsonDefinition = fmt.Sprintf("`json:\"%s\"`", key)
			} else {
				jsonDefinition = fmt.Sprintf("`json:\"%s,omitempty\"`", key)
				required = false
			}
			propertyType, err := getPropertyTypeAndCreateNewBlock(key, value, filePath, required)
			if err != nil {
				return err
			}
			statusProperties = append(statusProperties, blockProperty{
				Name: strings.Title(key),
				Type: propertyType,
				JSON: jsonDefinition,
			})
		}
		sort.Slice(statusProperties, func(i, j int) bool { return statusProperties[i].Name < statusProperties[j].Name })
		if err := statusTemplate.Execute(file, &templateOptionsSpecStatusBlock{
			Kind:       crd.Spec.Names.Kind,
			Properties: statusProperties,
		}); err != nil {
			return err
		}
	}
	//additional blocks
	for {
		for propertyName, propertyValue := range blockProperties {
			var chieldProperties []blockProperty
			blockTemplate, err := template.New("block-tmpl").Parse(templates.TemplateBlock)
			if err != nil {
				return err
			}
			for key, value := range propertyValue.Properties {
				var jsonDefinition string
				required := true
				if contains(propertyValue.Required, key) {
					jsonDefinition = fmt.Sprintf("`json:\"%s\"`", key)
				} else {
					jsonDefinition = fmt.Sprintf("`json:\"%s,omitempty\"`", key)
					required = false
				}
				propertyType, err := getPropertyTypeAndCreateNewBlock(key, value, filePath, required)
				if err != nil {
					return err
				}
				chieldProperties = append(chieldProperties, blockProperty{
					Name: strings.Title(key),
					Type: propertyType,
					JSON: jsonDefinition,
				})
			}

			sort.Slice(chieldProperties, func(i, j int) bool { return chieldProperties[i].Name < chieldProperties[j].Name })
			re := regexp.MustCompile(`[\n\r]+`)
			description := re.ReplaceAll([]byte(propertyValue.Description), []byte("\n// "))

			if err := blockTemplate.Execute(file, &templateOptionsBlock{
				Kind:        propertyName,
				Properties:  chieldProperties,
				Description: string(description),
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

func getPropertyTypeAndCreateNewBlock(key string, value OpenAPIV3Schema, filePath string, required bool) (string, error) {
	switch value.Type {
	case "array":
		kind := trimSuffix(strings.Title(key), "s")
		if value.Items.Type == "object" {
			exists, err := propertyBlockExists(filePath, kind)
			if err != nil {
				return "", err
			}
			if !exists {
				blockProperties[kind] = value.Items
			}
			if required {
				return fmt.Sprintf("[]%s", kind), nil
			}
			return fmt.Sprintf("[]*%s", kind), nil
		}
		if required {
			return fmt.Sprintf("[]%s", value.Items.Type), nil
		}
		return fmt.Sprintf("[]*%s", value.Items.Type), nil

	case "object":
		kind := strings.Title(key)
		if value.AdditionalProperties != nil {
			return fmt.Sprintf("map[string]%s", value.AdditionalProperties.Type), nil
		}
		exists, err := propertyBlockExists(filePath, kind)
		if err != nil {
			return "", err
		}
		if !exists {
			blockProperties[strings.Title(key)] = &value
		}
		if required {
			return kind, nil
		}
		return fmt.Sprintf("*%s", kind), nil
	case "boolean":
		if required {
			return "bool", nil
		}
		return "*bool", nil
	case "integer":
		if required {
			return "int", nil
		}
		return "*int", nil
	case "string":
		return "string", nil
	default:
		if required {
			return value.Type, nil
		}
		return fmt.Sprintf("*%s", value.Type), nil
	}
}

func propertyBlockExists(filePath string, kind string) (bool, error) {
	folderPath := filepath.Dir(filePath)
	blockExists := false
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return blockExists, err
	}
	for _, file := range files {
		path := filepath.Join(folderPath, file.Name())
		read, err := ioutil.ReadFile(path)
		if err != nil {
			return blockExists, err
		}
		if strings.Contains(string(read), fmt.Sprintf("type %s struct {", kind)) {
			blockExists = true
		}
	}
	return blockExists, nil
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

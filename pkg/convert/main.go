package convert

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

type groupVersion struct {
	Group   string
	Version string
}

// CRD2Api includes the convertion logic
func CRD2Api(crdFolder string, outputFolder string, licenseFilePath string) error {
	var apiGroupVersions []groupVersion
	var crds []*CustomResourceDefinition

	// Reading license
	license := ""
	if licenseFilePath != "" {
		licenseByte, err := ioutil.ReadFile(licenseFilePath)
		if err != nil {
			return err
		}
		license = string(licenseByte)
	}

	// Reading CRD files
	var files []string
	if err := filepath.Walk(crdFolder, func(path string, info os.FileInfo, err error) error {
		if path != crdFolder {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return err
	}

	for _, file := range files {
		log.Printf("Reading file: %s", file)
		crd, err := readCrdFile(file)
		crds = append(crds, crd)
		if err != nil {
			return err
		}
		gv := groupVersion{
			Group:   crd.Spec.Group,
			Version: crd.Spec.Version,
		}
		foundGroupVersion := false
		for _, v := range apiGroupVersions {
			if v.Group == gv.Group && v.Version == gv.Version {
				foundGroupVersion = true
			}
		}
		if foundGroupVersion == false {
			apiGroupVersions = append(apiGroupVersions, gv)
		}
	}

	// Generate Addtoscheme files and folders
	for _, gv := range apiGroupVersions {
		groupFolderPath := path.Join(outputFolder, gv.Group)
		if _, err := os.Stat(groupFolderPath); os.IsNotExist(err) {
			if err := os.Mkdir(groupFolderPath, 0766); err != nil {
				return err
			}
		}
		versionFolderPath := path.Join(groupFolderPath, gv.Version)
		if _, err := os.Stat(versionFolderPath); os.IsNotExist(err) {
			if err := os.Mkdir(versionFolderPath, 0766); err != nil {
				return err
			}
		}
		if err := createGroupversionInfoGoFile(path.Join(versionFolderPath, "groupversion_info.go"), &templateOptionsGroupversionInfo{
			Group:   gv.Group,
			Version: gv.Version,
			License: license,
		}); err != nil {
			return err
		}
	}

	// Generate type files
	for _, crd := range crds {
		filepath := path.Join(outputFolder, crd.Spec.Group, crd.Spec.Version, fmt.Sprintf("%s_types.go", crd.Spec.Names.Singular))
		if err := createTypesGoFile(filepath, crd, license); err != nil {
			return err
		}
	}
	return nil
}

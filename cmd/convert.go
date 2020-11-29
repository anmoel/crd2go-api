/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"

	"github.com/anmoel/crd2go-api/pkg/convert"
	"github.com/spf13/cobra"
)

var licenseFilePath string

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert <CRD-FOLDER> <OUTPUT-FOLDER>",
	Short: "Convert kubernetes CRDs to go api definitions",
	Long:  `Convert all kubernetes CustonResourceDefinitions of the import-folder to a golang api definition, stored in the output folder.`,
	Args:  cobra.MinimumNArgs(2),

	Run: func(cmd *cobra.Command, args []string) {
		if err := convert.CRD2Api(args[0], args[1], licenseFilePath); err != nil {
			log.Panicf("Could not generate api definitions. Error %v", err)
		} else {
			log.Println("Api definitions generated")
		}
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.
	convertCmd.PersistentFlags().StringVarP(&licenseFilePath, "license-file-path", "l", "", "The path to the license file is used by the templates")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// convertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

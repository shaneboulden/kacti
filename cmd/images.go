/*
Copyright © 2023 Shane Boulden (sboulden@redhat.com)

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
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var listImg = false

var cve2image = make(map[string]string)

func initCVEs() {
	// Populate the map with sample data
	cve2image["CVE-2021-44228"] = "quay.io/kacti/log4shell@sha256:f72efa1cb3533220212bc49716a4693b448106b84ca259c20422ab387972eed9"
}

var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "Verify and list Kacti images",
	Long:  `Verify and list Kacti images`,
	Run: func(cmd *cobra.Command, args []string) {
		if listImg {
			listImages()
		}
	},
}

func listImages() {
	for cve := range cve2image {
		cveImage := cve2image[cve]
		fmt.Print(color.YellowString("CVE: %s -> Image: %s\n", cve, cveImage))
	}
}

func init() {
	imagesCmd.Flags().BoolVarP(&listImg, "list", "l", false, "List supported Kacti images")
	rootCmd.AddCommand(imagesCmd)
}

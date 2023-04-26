/*
Copyright © 2023 Shane Boulden (shane.boulden@gmail.com)

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
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Test struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	Image     string `yaml:image"`
}

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Funcionally test admission control",
	Long: `Perform functional verification tests against Kubernetes admission controllers.

Tests are specififed in files that reference test names, namespaces and images. For example:

kacti-tests:
- name: pwnkit
  image: quay.io/the-worst-containers/pwnkit:v0.2
  namespace: app-deploy`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runTest(strings.Join(args, " "))
	},
}

func runTest(file string) {
	fmt.Println("Using tests from: " + file)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	var testConfig Test
	err = yaml.Unmarshal(data, &testConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the test config
	fmt.Println(("Running test: " + testConfig.Name + " { ns: " + testConfig.Namespace + " // img:" + testConfig.Image + " }"))
}

func init() {
	rootCmd.AddCommand(testCmd)
}

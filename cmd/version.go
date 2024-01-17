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

	"github.com/spf13/cobra"
)

var (
	// set with ldflags by slsa-golangreleaser
	Version    string
	Commit     string
	CommitDate string
	TreeState  string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\nCommit: %s\nCommit Date: %s\nTree State: %s\n", Version, Commit, CommitDate, TreeState)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// Copyright 息 2016 Shigeyuki Fujishima <shigeyuki.fujishima@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/dafiti-group/k8s-values-updater/pkg/bump"
	"github.com/spf13/cobra"
)

var b bump.Bump

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "bump",
	Short: "Bump value on file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		b.DryRun = dryRun
		if err = b.Run(); err != nil {
			fmt.Println("Error")
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringVar(&b.RemoteName, "remote-name", "origin", "")
	addCmd.PersistentFlags().StringVar(&b.DirPath, "dir-path", "deploy/*", "File Path")
	addCmd.PersistentFlags().StringVar(&b.FileNames, "file-names", "values.yaml", "File Path")
	addCmd.PersistentFlags().StringVar(&b.ReplaceWith, "replace-with", "", "If passed will try to merge this value with the values yaml")
	addCmd.PersistentFlags().StringVar(&b.ChartName, "chart-name", "", "The name of the subchart")
	addCmd.PersistentFlags().BoolVarP(&b.IsRoot, "is-root", "", false, "If set will define that the values to be changed has no subchart")
	addCmd.PersistentFlags().StringVar(&b.Email, "email", "k8s-values-updater@mailinator.com", "Email that will commit")
	addCmd.PersistentFlags().StringVar(&b.PrID, "pr-id", "", "Pull Request ID")
	addCmd.PersistentFlags().StringVar(&b.Tag, "tag", "", "Image Tag")
	addCmd.PersistentFlags().StringVar(&b.RefSpec, "refspec", "", "refspec")
	addCmd.PersistentFlags().StringVar(&b.Branch, "branch", "", "branch (required)")
}

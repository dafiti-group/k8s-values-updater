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
	"os"

	"github.com/dafiti-group/k8s-values-updater/pkg/bump"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var b = bump.New(logrus.New())
var githubAccesToken string
var separator = ","

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "bump",
	Short: "Bump value on file",
	Long:  ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		// Get global var
		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			logrus.Panic(os.Stderr, err)
			os.Exit(2)
		}

		// Init package
		err = b.Init(githubAccesToken, dryRun, separator)
		if err != nil {
			logrus.Panic(os.Stderr, err)
			os.Exit(2)
		}

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		//
		if err := b.Run(); err != nil {
			logrus.Panic(os.Stderr, err)
			os.Exit(2)
		}
	},
}

func init() {
	cobra.OnInitialize(func() {
		initEnvs(b)
		v.AutomaticEnv()
	})

	addCmd.PersistentFlags().BoolVarP(&b.IsRoot, "is-root", "", false, "If set will define that the values to be changed has no subchart")
	addCmd.PersistentFlags().StringVar(&b.ChartName, "chart-name", "", "The name of the subchart")
	addCmd.PersistentFlags().StringVar(&b.DirPath, "dir-path", "deploy/*", "File Path")
	addCmd.PersistentFlags().StringVar(&b.FileNames, "file-names", "values.yaml", "File Path")
	addCmd.PersistentFlags().StringVar(&b.RemoteName, "remote-name", "origin", "")
	addCmd.PersistentFlags().StringVar(&b.ReplaceWith, "replace-with", "", "If passed will try to merge this value with the values yaml")
	addCmd.PersistentFlags().StringVar(&githubAccesToken, "github-access-token", "", "Github Acccess Token (Will use if set $GITHUB_ACCESS_TOKEN)")
	addCmd.PersistentFlags().StringVarP(&b.Branch, "branch", "b", "", "Branch (Will use if set $CIRCLE_BRANCH)")
	addCmd.PersistentFlags().StringVarP(&b.Email, "email", "e", "k8s-values-updater@mailinator.com", "Email that will commit")
	addCmd.PersistentFlags().StringVarP(&b.PrID, "pr-id", "p", "", "Pull Request ID")
	addCmd.PersistentFlags().StringVarP(&b.Remote.Org, "org", "o", "", "Organization (Will use if set $CIRCLE_PROJECT_USERNAME)")
	addCmd.PersistentFlags().StringVarP(&b.Remote.Repo, "repo", "r", "", "Repository (Will use if set $CIRCLE_PROJECT_REPONAME)")
	addCmd.PersistentFlags().StringVarP(&b.Tag, "tag", "t", "", "Image Tag")
	addCmd.PersistentFlags().StringVarP(&b.WorkDir, "workdir", "w", ".", "Workdir")

	RootCmd.AddCommand(addCmd)
}

func initEnvs(b *bump.Bump) {
	v.BindEnv("github_access_token")

	v.SetEnvPrefix("circle")
	v.BindEnv("branch")
	v.BindEnv("project_reponame")
	v.BindEnv("project_username")

	// TODO: try to autoload this
	githubAccesToken = v.GetString("github_access_token")
	b.Remote.Repo = v.GetString("project_reponame")
	b.Remote.Org = v.GetString("project_username")
	b.Branch = v.GetString("branch")

	// Take the envs and load it on the struct
	err := v.Unmarshal(b)
	if err != nil {
		logrus.Panic(os.Stderr, err)
		panic(-1)
	}
}

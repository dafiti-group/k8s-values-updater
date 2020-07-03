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
	"github.com/sirupsen/logrus"

	"github.com/dafiti-group/k8s-values-updater/pkg/bump/file"
	"github.com/dafiti-group/k8s-values-updater/pkg/bump/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var b bump.Bump
var g git.Git
var f file.File
var githubAccesToken string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "bump",
	Short: "Bump value on file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		//
		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			logrus.Panic(os.Stderr, err)
			os.Exit(2)
		}

		//
		err = b.Init(&g, &f, githubAccesToken, dryRun, logrus.New())
		if err != nil {
			logrus.Panic(os.Stderr, err)
			os.Exit(2)
		}

		//
		if err = b.Run(); err != nil {
			logrus.Panic(os.Stderr, err)
			os.Exit(2)
		}
	},
}

func init() {
	cobra.OnInitialize(func() {
		initBump(&b)
	})

	RootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().BoolVarP(&f.IsRoot, "is-root", "", false, "If set will define that the values to be changed has no subchart")
	addCmd.PersistentFlags().StringVar(&b.DirPath, "dir-path", "deploy/*", "File Path")
	addCmd.PersistentFlags().StringVar(&b.FileNames, "file-names", "values.yaml", "File Path")
	addCmd.PersistentFlags().StringVar(&githubAccesToken, "github-access-token", "", "Github Acccess Token")
	addCmd.PersistentFlags().StringVar(&f.ChartName, "chart-name", "", "The name of the subchart")
	addCmd.PersistentFlags().StringVar(&f.ReplaceWith, "replace-with", "", "If passed will try to merge this value with the values yaml")
	addCmd.PersistentFlags().StringVar(&g.RemoteName, "remote-name", "origin", "")
	addCmd.PersistentFlags().StringVarP(&f.PrID, "pr-id", "p", "", "Pull Request ID")
	addCmd.PersistentFlags().StringVarP(&f.Tag, "tag", "t", "", "Image Tag")
	addCmd.PersistentFlags().StringVarP(&g.Email, "email", "e", "k8s-values-updater@mailinator.com", "Email that will commit")
	addCmd.PersistentFlags().StringVarP(&g.WorkDir, "workdir", "w", ".", "Workdir")
}

func initBump(b *bump.Bump) {
	v := viper.New()
	// TODO: Maybe set some default envs
	v.BindEnv("GITHUB_ACCESS_TOKEN")

	// Will not persist this on a struct
	githubAccesToken = v.GetString("GITHUB_ACCESS_TOKEN")

	// Take the envs and load it on the struct
	err := v.Unmarshal(b)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

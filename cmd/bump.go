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
	"github.com/k0kubun/pp"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var b = bump.New(logrus.New())
var githubAccesToken string
var separator = ","

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "bump",
	Short: "Bump value on file",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		pp.Println(args)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		//
		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			logrus.Panic(os.Stderr, err)
			os.Exit(2)
		}

		pp.Println(b.URL)
		pp.Println(b.Branch)
		pp.Println(githubAccesToken)
		panic("---")

		//
		err = b.Init(githubAccesToken, dryRun, separator)
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
		viper.AutomaticEnv()
		initEnvs(b)
		addCmd.Flags().VisitAll(func(f *pflag.Flag) {
			if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
				addCmd.Flags().Set(f.Name, viper.GetString(f.Name))
			}
		})
	})

	RootCmd.AddCommand(addCmd)

	addCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) (err error) {
		initEnvs(b)
		return nil
	}

	addCmd.PersistentFlags().BoolVarP(&b.IsRoot, "is-root", "", false, "If set will define that the values to be changed has no subchart")
	addCmd.PersistentFlags().StringVar(&b.DirPath, "dir-path", "deploy/*", "File Path")
	addCmd.PersistentFlags().StringVar(&b.FileNames, "file-names", "values.yaml", "File Path")
	addCmd.PersistentFlags().StringVar(&githubAccesToken, "github-access-token", "", "Github Acccess Token")
	addCmd.PersistentFlags().StringVar(&b.ChartName, "chart-name", "", "The name of the subchart")
	addCmd.PersistentFlags().StringVar(&b.ReplaceWith, "replace-with", "", "If passed will try to merge this value with the values yaml")
	addCmd.PersistentFlags().StringVar(&b.RemoteName, "remote-name", "origin", "")
	addCmd.PersistentFlags().StringVarP(&b.PrID, "pr-id", "p", "", "Pull Request ID")
	addCmd.PersistentFlags().StringVarP(&b.Tag, "tag", "t", "", "Image Tag")
	addCmd.PersistentFlags().StringVarP(&b.Email, "email", "e", "k8s-values-updater@mailinator.com", "Email that will commit")
	addCmd.PersistentFlags().StringVarP(&b.WorkDir, "workdir", "w", ".", "Workdir")

	// addCmd.MarkPersistentFlagRequired("branch")
	// addCmd.MarkPersistentFlagRequired("github-access-token")
	// addCmd.MarkPersistentFlagRequired("url")
}

func initEnvs(b *bump.Bump) {
	v := viper.New()
	// TODO: Maybe set some default envs
	v.SetEnvPrefix("circle")
	v.BindEnv("github_access_token")
	v.BindEnv("project_username")
	v.BindEnv("branch")
	v.BindEnv("project_reponame")
	// CIRCLE_PROJECT_USERNAME=dafiti-group
	// CIRCLE_PROJECT_REPONAME=k8s-values-updater
	// CIRCLE_BRANCH=feature/auth-only-with-https

	// Will not persist this on a struct
	b.SetAuth(v.GetString("GITHUB_ACCESS_TOKEN"))
	// branch = v.GetString("CIRCLE_BRANCH")
	// org = v.GetString("CIRCLE_PROJECT_USERNAME")
	// repo = v.GetString("CIRCLE_PROJECT_REPONAME")

	// Take the envs and load it on the struct
	err := v.Unmarshal(b)
	if err != nil {
		logrus.Panic(os.Stderr, err)
		os.Exit(2)
		panic(2)
	}
}

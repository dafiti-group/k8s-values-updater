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
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var logLevel string
var verbose bool
var dryRun bool
var v = viper.New()

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "k8s-values-updater",
	Short: "Updates the k8s values",
	Long:  ``,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logrus.Fatal(err)
		os.Exit(-1)
	}
}

func init() {
	// Init config, will not be affected by log level change
	cobra.OnInitialize(initConfig)

	// Set Log level
	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := setUpLogs(os.Stdout, verbose, logLevel); err != nil {
			return err
		}

		logrus.WithFields(logrus.Fields{
			"level": logrus.GetLevel(),
		}).Info("start")

		return nil
	}

	// Flags
	RootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Dry Run")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Dry Run")
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config.yaml)")
	RootCmd.PersistentFlags().StringVar(&logLevel, "log-level", logrus.WarnLevel.String(), "Log level (debug, info, warn, error, fatal, panic")
}

//setUpLogs set the log output ans the log level
func setUpLogs(out io.Writer, verbose bool, level string) error {
	// Log Level takes presedence
	if level == logrus.WarnLevel.String() && verbose {
		level = logrus.DebugLevel.String()
	}
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	return nil
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	logrus.Info("read config")

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logrus.Warn(err)
	}

	v.SetConfigName(".config") // name of config file (without extension)
	if cfgFile != "" {         // enable ability to specify config file via flag
		logrus.Info(">>> cfgFile: ", cfgFile)
		viper.SetConfigFile(cfgFile)
		configDir := path.Dir(cfgFile)
		if configDir != "." && configDir != dir {
			viper.AddConfigPath(configDir)
		}
	}

	v.AddConfigPath(dir)
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME")
	v.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := v.ReadInConfig(); err == nil {
		logrus.Info("Using config file:", v.ConfigFileUsed())
	} else {
		logrus.Error(err)
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		logrus.Info("Config file changed:", e.Name)
	})
}

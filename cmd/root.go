/*
Copyright © 2022 Ismet Handzic

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
	"os"
	"runtime"

	"github.com/fatih/color"
	"github.com/ismet55555/the-mouse-is-lava/lava"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "the-mouse-is-lava",
	Short: "A little program that helps you keep your fingers off the mouse",
	Long:  `==========================================
    The
       mouse
            is
               ╦  ┌─┐┬  ┬┌─┐
               ║  ├─┤└┐┌┘├─┤
               ╩═╝┴ ┴ └┘ ┴ ┴
==========================================
This stupid program helps you keep your fingers off the mouse`,
	Run: func(cmd *cobra.Command, args []string) {
		l := lava.Lava{}
		l.Start()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Initial setup
// Here you will define your flags and configuration settings.
// Cobra supports persistent flags, which, if defined here,
// will be global for your application.
func init() {
	// Check for windows platform
	if runtime.GOOS == "windows" {
		color.HiYellow("Sorry. Windows is not fully supported yet. Please check back with later version")
		os.Exit(1)
	}

	cobra.OnInitialize(initConfig)

    // Persistent CLI flags - Will be global for application
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.the-mouse-is-lava.yaml)")
    viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose debug logs")
    viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))


	// Local CLI flags - Will only run when this action is called directly not all commands
	rootCmd.Flags().Bool("no-systray", false, "Disable system systray menu")
    viper.BindPFlag("noSystray", rootCmd.Flags().Lookup("no-systray"))

	rootCmd.Flags().Bool("no-intro", false, "Disable cool intro animation")
    viper.BindPFlag("noIntro", rootCmd.Flags().Lookup("no-intro"))

	rootCmd.Flags().Bool("detach", false, "Detach process to run in background")
    viper.BindPFlag("detach", rootCmd.Flags().Lookup("detach"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Configure logger
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
    if viper.GetBool("verbose") {
        log.SetLevel(log.DebugLevel)
    } else {
        log.SetLevel(log.InfoLevel)
    }

    // Default configuration values
    viper.SetDefault("initAnimation", true)
    viper.SetDefault("initGracePeriod", 2)
    viper.SetDefault("initPause", true)
    viper.SetDefault("gracePeriodDuration", 3)
    viper.SetDefault("gracePeriod", false)
    viper.SetDefault("sensitivity", 8.0)
    viper.SetDefault("enableSystray", true)

    if cfgFile != "" {
        // Use config file from the flag.
        viper.SetConfigFile(cfgFile)
    } else {
        // Find home directory.
        home, err := os.UserHomeDir()
        cobra.CheckErr(err)

        // Search config in home directory with name ".the-mouse-is-lava" (without extension).
        viper.AddConfigPath(home)
        viper.AddConfigPath(".")
        viper.SetConfigType("yaml")
        viper.SetConfigName(".the-mouse-is-lava")
    }

    // Read in environment variables that match prefix
    viper.AutomaticEnv()

    // If a config file is found, read it in.
    if err := viper.ReadInConfig(); err == nil {
        log.Debug("Using config file:", viper.ConfigFileUsed())
    }
}

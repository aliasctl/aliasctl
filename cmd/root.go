package cmd

import (
	"fmt"
	"os"

	"github.com/aliasctl/aliasctl/internal/config"
	"github.com/spf13/cobra"
)

var cfg *config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	// ...existing code...
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
// If there is an error during execution, the program will exit with status code 1.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// ...existing code...
}

// initConfig reads in config file and ENV variables if set.
// It loads the application configuration using the config package
// and stores it in the global cfg variable for use by other commands.
// If configuration fails to load, the program terminates with an error message.
func initConfig() {
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}
}

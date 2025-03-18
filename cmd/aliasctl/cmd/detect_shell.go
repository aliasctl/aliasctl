package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// detectShellCmd represents the detect-shell command
var detectShellCmd = &cobra.Command{
	Use:   "detect-shell",
	Short: "Show detected shell and alias file",
	Long:  `Display information about the detected shell type and alias file path.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Detected shell: %s\n", am.Shell)
		fmt.Printf("Alias file: %s\n", am.AliasFile)
		fmt.Printf("Config directory: %s\n", am.ConfigDir)

		// Add additional helpful information
		if _, err := os.Stat(am.AliasFile); os.IsNotExist(err) {
			fmt.Printf("\nNote: The alias file does not exist yet. It will be created when you add aliases.\n")
		}

		if _, err := os.Stat(am.ConfigDir); os.IsNotExist(err) {
			fmt.Printf("\nNote: The config directory doesn't exist yet. It will be created automatically.\n")
		}

		fmt.Printf("\nTo change shell type: aliasctl set-shell <shell-type>\n")
		fmt.Printf("To change alias file: aliasctl set-file <file-path>\n")
	},
}

func init() {
	rootCmd.AddCommand(detectShellCmd)
}

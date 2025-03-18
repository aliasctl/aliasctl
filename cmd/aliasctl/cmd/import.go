package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// importCmd represents the import command which loads aliases from an existing shell configuration file.
// It reads the current shell's format and extracts any alias definitions it can find.
// The command will validate that the file exists before attempting to import.
// Example usage: aliasctl import
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import aliases from shell configuration",
	Long:  `Import aliases from the current shell's configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if the alias file exists first
		if _, err := os.Stat(am.AliasFile); os.IsNotExist(err) {
			return fmt.Errorf("shell configuration file not found at '%s'\n\nUse 'aliasctl set-file' to specify a different location or create the file manually", am.AliasFile)
		}

		if err := am.ImportAliasesFromShell(); err != nil {
			return fmt.Errorf("failed to import aliases from shell configuration: %w\n\nMake sure '%s' contains valid alias definitions for %s shell", err, am.AliasFile, am.Shell)
		}

		fmt.Println("Aliases successfully imported from shell configuration")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}

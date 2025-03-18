package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

// setFileCmd represents the set-file command which configures the path where aliases will be stored.
// This allows users to choose a custom location for their shell configuration file.
// The path can be absolute or relative; if relative, it will be resolved to an absolute path.
// Example usage: aliasctl set-file ~/.my_aliases
var setFileCmd = &cobra.Command{
	Use:   "set-file [alias-file-path]",
	Short: "Manually set the alias file path",
	Long:  `Manually set the path to the file where aliases will be stored.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		absPath, _ := filepath.Abs(filePath)

		if err := am.SetAliasFile(filePath); err != nil {
			return fmt.Errorf("failed to set alias file path to '%s': %w\n\nEnsure the directory exists and you have write permissions", absPath, err)
		}

		fmt.Printf("Alias file successfully set to %s\n", am.AliasFile)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setFileCmd)
}

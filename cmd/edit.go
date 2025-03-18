package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [alias]",
	Short: "Edit an alias file",
	Long:  `Edit an existing alias file using your default editor.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		aliasName := args[0]
		aliasPath := filepath.Join(cfg.AliasesDir, aliasName)

		// Check if the alias file exists
		if _, err := os.Stat(aliasPath); os.IsNotExist(err) {
			return fmt.Errorf("alias '%s' does not exist", aliasName)
		}

		// Open the file with the default editor
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = cfg.DefaultEditor
		}

		editorCmd := exec.Command(editor, aliasPath)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		return editorCmd.Run()
	},
}

// init registers the edit command and its flags to the root command.
// It's called automatically when the package is initialized.
func init() {
	rootCmd.AddCommand(editCmd)
}

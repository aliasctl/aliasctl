package cmd

import (
	"fmt"

	"github.com/aliasctl/aliasctl/pkg/aliasctl"
	"github.com/spf13/cobra"
)

// encryptAPIKeysCmd represents the encrypt-api-keys command which secures API keys using encryption.
// This command encrypts any plaintext API keys in the configuration and stores the encrypted
// version instead. The encryption key is stored separately for security.
// Example usage: aliasctl encrypt-api-keys
var encryptAPIKeysCmd = &cobra.Command{
	Use:   "encrypt-api-keys",
	Short: "Encrypt API keys in configuration",
	Long:  `Encrypt API keys stored in the configuration file for security.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := am.EncryptAPIKeys(); err != nil {
			return fmt.Errorf("failed to encrypt API keys: %w\n\nEnsure you have write permissions to %s and the directory exists", err, am.EncryptionKey)
		}
		fmt.Println("API keys successfully encrypted")
		fmt.Printf("Encryption key stored at: %s\n", am.EncryptionKey)
		fmt.Println("WARNING: Keep this key file secure as it's needed to decrypt your API keys.")
		fmt.Println("Consider backing up this key file to a secure location.")
		return nil
	},
}

// disableEncryptionCmd represents the disable-encryption command which reverts to plaintext API keys.
// This command decrypts any encrypted API keys in the configuration and stores them in plaintext.
// This is not recommended for security reasons but may be necessary in some cases.
// Example usage: aliasctl disable-encryption
var disableEncryptionCmd = &cobra.Command{
	Use:   "disable-encryption",
	Short: "Disable API key encryption",
	Long:  `Disable encryption and revert to plaintext API keys (not recommended).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := am.DisableEncryption(); err != nil {
			if _, ok := err.(*aliasctl.KeyFileNotFoundError); ok {
				return fmt.Errorf("encryption key not found at %s: %w\n\nIf you've lost your encryption key, you'll need to reconfigure your API providers", am.EncryptionKey, err)
			}
			return fmt.Errorf("failed to disable encryption: %w", err)
		}
		fmt.Println("API key encryption successfully disabled")
		fmt.Println("WARNING: API keys are now stored in plaintext. This is not recommended for security reasons.")
		fmt.Println("Your API keys are now vulnerable if someone gains access to your configuration file.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(encryptAPIKeysCmd)
	rootCmd.AddCommand(disableEncryptionCmd)
}

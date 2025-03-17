// New file that handles configuration paths
package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetConfigDir returns the configuration directory for the application.
func GetConfigDir() string {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			homeDir, _ := os.UserHomeDir()
			appData = filepath.Join(homeDir, "AppData", "Roaming")
		}
		configDir = filepath.Join(appData, "AliasCtl")
	default:
		homeDir, _ := os.UserHomeDir()
		configDir = filepath.Join(homeDir, ".config", "aliasctl")
	}

	return configDir
}

// GetEncryptionKeyPath returns the path to the encryption key file
func GetEncryptionKeyPath(configDir string) string {
	return filepath.Join(configDir, ".keyring", "api_keys.key")
}

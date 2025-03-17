package aliasctl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// NewAliasManager creates a new AliasManager.
func NewAliasManager() *AliasManager {
	platform := runtime.GOOS
	configDir := getConfigDir()

	am := &AliasManager{
		Platform:       platform,
		Aliases:        make(map[string]AliasCommands),
		AIConfigured:   false,
		AIProviders:    make(map[string]AIProvider),
		ConfigDir:      configDir,
		AliasStore:     filepath.Join(configDir, "aliases.json"),
		ConfigFile:     filepath.Join(configDir, "config.json"),
		EncryptionKey:  GetEncryptionKeyPath(configDir),
		EncryptionUsed: false,
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("Warning: couldn't create config directory: %v\n", err)
	}

	if err := am.LoadConfig(); err != nil {
		shell, aliasFile := DetectShellAndAliasFile(platform)
		am.Shell = shell
		am.AliasFile = aliasFile
		am.SaveConfig()
	}

	return am
}

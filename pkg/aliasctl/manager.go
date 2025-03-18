package aliasctl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/aliasctl/aliasctl/pkg/aliasctl/ai"
)

// NewAliasManager creates a new AliasManager.
func NewAliasManager() *AliasManager {
	platform := runtime.GOOS
	configDir := getConfigDir()

	// Fix the GetEncryptionKeyPath call to handle both return values
	encryptionKeyPath, err := GetEncryptionKeyPath(configDir)
	if err != nil {
		fmt.Printf("Warning: failed to get encryption key path: %v\n", err)
		encryptionKeyPath = filepath.Join(configDir, "encryption.key") // Fallback path
	}

	am := &AliasManager{
		Platform:       platform,
		Aliases:        make(map[string]AliasCommands),
		AIConfigured:   false,
		aiManager:      ai.NewManager(),
		ConfigDir:      configDir,
		AliasStore:     filepath.Join(configDir, "aliases.json"),
		ConfigFile:     filepath.Join(configDir, "config.json"),
		EncryptionKey:  encryptionKeyPath,
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

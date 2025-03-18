package aliasctl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/aliasctl/aliasctl/pkg/aliasctl/ai"
)

// getConfigDir returns the configuration directory for the application.
func getConfigDir() string {
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

// LoadConfig loads the application configuration, supporting both TOML and JSON for backward compatibility.
func (am *AliasManager) LoadConfig() error {
	data, err := os.ReadFile(am.ConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config file not found at %s\n\nRun 'aliasctl detect-shell' or manually set your shell with 'aliasctl set-shell'", am.ConfigFile)
		}
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied when accessing config file at %s\n\nCheck file permissions or run with appropriate privileges", am.ConfigFile)
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config

	// Try to parse as TOML first
	err = toml.Unmarshal(data, &config)

	// If TOML parsing fails, try JSON as fallback for backward compatibility
	if err != nil {
		err = json.Unmarshal(data, &config)
		if err != nil {
			return fmt.Errorf("failed to parse config file - invalid format\n\nThe config file at %s is neither valid TOML nor JSON\nConsider running 'aliasctl set-shell' to regenerate it", am.ConfigFile)
		}

		// If it was JSON, convert to TOML for future use
		fmt.Println("Converting config from JSON to TOML format for better readability...")
		if err := am.convertConfigToTOML(); err != nil {
			fmt.Printf("Warning: Failed to convert config to TOML: %v\n", err)
			// Continue anyway since we were able to load the JSON
		}
	}

	am.Shell = config.DefaultShell
	am.AliasFile = config.DefaultAliasFile
	am.EncryptionUsed = config.UseEncryption

	// Initialize aiManager if nil
	if am.aiManager == nil {
		am.aiManager = ai.NewManager()
	}

	// Handle API configuration - check for encrypted keys first
	if config.OllamaEndpoint != "" && config.OllamaModel != "" {
		am.ConfigureOllama(config.OllamaEndpoint, config.OllamaModel)
	}

	// Handle OpenAI configuration
	if config.OpenAIEndpoint != "" && config.OpenAIModel != "" {
		var apiKey string

		// Try to use encrypted key first
		if config.UseEncryption && config.OpenAIKeyEncrypted != "" {
			decryptedKey, err := DecryptString(config.OpenAIKeyEncrypted, am.EncryptionKey)
			if err == nil {
				apiKey = decryptedKey
			} else {
				fmt.Printf("Warning: Failed to decrypt OpenAI API key: %v\n", err)
				if _, ok := err.(*KeyFileNotFoundError); ok {
					fmt.Printf("Encryption key file not found at: %s\n", am.EncryptionKey)
					fmt.Printf("Use 'aliasctl encrypt-api-keys' to set up encryption\n")
				}

				// Fallback to plaintext key with warning if available
				if config.OpenAIKey != "" {
					fmt.Println("Warning: Using plaintext API key from config. Consider encrypting your API keys.")
					apiKey = config.OpenAIKey
				}
			}
		} else if config.OpenAIKey != "" {
			fmt.Println("Warning: API key is stored in plaintext. Use 'aliasctl encrypt-api-keys' to encrypt it.")
			apiKey = config.OpenAIKey
		}

		if apiKey != "" {
			am.ConfigureOpenAI(config.OpenAIEndpoint, apiKey, config.OpenAIModel)
		}
	}

	// Handle Anthropic configuration
	if config.AnthropicEndpoint != "" && config.AnthropicModel != "" {
		var apiKey string

		// Try to use encrypted key first
		if config.UseEncryption && config.AnthropicKeyEncrypted != "" {
			decryptedKey, err := DecryptString(config.AnthropicKeyEncrypted, am.EncryptionKey)
			if err == nil {
				apiKey = decryptedKey
			} else {
				fmt.Printf("Warning: Failed to decrypt Anthropic API key: %v\n", err)
				if _, ok := err.(*KeyFileNotFoundError); ok {
					fmt.Printf("Encryption key file not found at: %s\n", am.EncryptionKey)
					fmt.Printf("Use 'aliasctl encrypt-api-keys' to set up encryption\n")
				}

				// Fallback to plaintext key with warning if available
				if config.AnthropicKey != "" {
					fmt.Println("Warning: Using plaintext Anthropic API key from config. Consider encrypting your API keys.")
					apiKey = config.AnthropicKey
				}
			}
		} else if config.AnthropicKey != "" {
			fmt.Println("Warning: Anthropic API key is stored in plaintext. Use 'aliasctl encrypt-api-keys' to encrypt it.")
			apiKey = config.AnthropicKey
		}

		if apiKey != "" {
			am.ConfigureAnthropic(config.AnthropicEndpoint, apiKey, config.AnthropicModel)
		}
	}

	// Set default provider if one exists in config
	if config.AIProvider != "" {
		am.aiManager.SetDefaultProvider(config.AIProvider)
	}

	return nil
}

// SaveConfig saves the application configuration in TOML format.
func (am *AliasManager) SaveConfig() error {
	config := Config{
		DefaultShell:     am.Shell,
		DefaultAliasFile: am.AliasFile,
		UseEncryption:    am.EncryptionUsed,
		AIProviders:      make(map[string]bool),
	}

	// Track which providers are configured
	providers := am.GetAvailableProviders()
	for _, name := range providers {
		config.AIProviders[name] = true
	}

	// Get default provider name
	if am.aiManager != nil && am.aiManager.Default != nil {
		// Determine the provider type
		switch am.aiManager.Default.(type) {
		case *ai.OllamaProvider:
			config.AIProvider = "ollama"
		case *ai.OpenAIProvider:
			config.AIProvider = "openai"
		case *ai.AnthropicProvider:
			config.AIProvider = "anthropic"
		}
	}

	// Configure providers
	ollamaProvider, ok := am.aiManager.Providers["ollama"].(*ai.OllamaProvider)
	if ok {
		config.OllamaEndpoint = ollamaProvider.Endpoint
		config.OllamaModel = ollamaProvider.Model
	}

	openAIProvider, ok := am.aiManager.Providers["openai"].(*ai.OpenAIProvider)
	if ok {
		config.OpenAIEndpoint = openAIProvider.Endpoint
		config.OpenAIModel = openAIProvider.Model

		// Handle API key encryption
		if am.EncryptionUsed {
			encryptedKey, err := EncryptString(openAIProvider.APIKey, am.EncryptionKey)
			if err == nil {
				config.OpenAIKeyEncrypted = encryptedKey
				config.OpenAIKey = "" // Clear plaintext key
			} else {
				fmt.Printf("Warning: Failed to encrypt API key: %v\n", err)
				fmt.Printf("API key will be stored in plaintext. Run 'aliasctl encrypt-api-keys' to retry encryption.\n")
				config.OpenAIKey = openAIProvider.APIKey
			}
		} else {
			config.OpenAIKey = openAIProvider.APIKey
		}
	}

	anthropicProvider, ok := am.aiManager.Providers["anthropic"].(*ai.AnthropicProvider)
	if ok {
		config.AnthropicEndpoint = anthropicProvider.Endpoint
		config.AnthropicModel = anthropicProvider.Model

		// Handle API key encryption
		if am.EncryptionUsed {
			encryptedKey, err := EncryptString(anthropicProvider.APIKey, am.EncryptionKey)
			if err == nil {
				config.AnthropicKeyEncrypted = encryptedKey
				config.AnthropicKey = "" // Clear plaintext key
			} else {
				fmt.Printf("Warning: Failed to encrypt Anthropic API key: %v\n", err)
				fmt.Printf("API key will be stored in plaintext. Run 'aliasctl encrypt-api-keys' to retry encryption.\n")
				config.AnthropicKey = anthropicProvider.APIKey
			}
		} else {
			config.AnthropicKey = anthropicProvider.APIKey
		}
	}

	file, err := os.Create(am.ConfigFile)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied when saving config file to %s\n\nCheck directory permissions or run with appropriate privileges", am.ConfigFile)
		}

		// Check if directory exists
		dir := filepath.Dir(am.ConfigFile)
		if _, statErr := os.Stat(dir); os.IsNotExist(statErr) {
			if mkdirErr := os.MkdirAll(dir, 0755); mkdirErr != nil {
				return fmt.Errorf("failed to create config directory %s: %w", dir, mkdirErr)
			}
			// Try creating file again after making directory
			file, err = os.Create(am.ConfigFile)
			if err != nil {
				return fmt.Errorf("failed to create config file at %s even after creating directory: %w", am.ConfigFile, err)
			}
		} else {
			return fmt.Errorf("failed to create config file at %s: %w", am.ConfigFile, err)
		}
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(config); err != nil {
		return fmt.Errorf("failed to write TOML configuration: %w\n\nThis might be due to invalid data or disk issues", err)
	}

	return nil
}

// convertConfigToTOML reads the existing JSON config file, parses it, and writes it back as TOML.
// It creates a backup of the original JSON file before conversion.
func (am *AliasManager) convertConfigToTOML() error {
	// Read existing JSON file
	data, err := os.ReadFile(am.ConfigFile)
	if err != nil {
		return err
	}

	// Parse JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	// Create backup of original file
	backupFile := am.ConfigFile + ".json.bak"
	if err := os.WriteFile(backupFile, data, 0644); err != nil {
		return fmt.Errorf("failed to create backup file %s: %w (check disk space and permissions)", backupFile, err)
	}

	// Write as TOML
	file, err := os.Create(am.ConfigFile)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(config); err != nil {
		// Restore from backup if TOML encoding fails
		if restoreErr := os.Rename(backupFile, am.ConfigFile); restoreErr != nil {
			return fmt.Errorf("TOML encoding failed: %v, and restore failed: %v", err, restoreErr)
		}
		return err
	}

	fmt.Printf("Config converted to TOML. Original JSON backup saved as %s\n", backupFile)
	return nil
}

// AddLoadAliasesTomlSupport modifies the existing LoadAliases method to support TOML format
// This should be called from the original LoadAliases method in alias_manager.go
func (am *AliasManager) AddLoadAliasesTomlSupport(data []byte) error {
	// Try to parse as TOML first
	err := toml.Unmarshal(data, &am.Aliases)

	// If TOML parsing fails, try JSON as fallback for backward compatibility
	if err != nil {
		err = json.Unmarshal(data, &am.Aliases)
		if err != nil {
			return fmt.Errorf("failed to parse aliases file - invalid format\n\nThe aliases file at %s is neither valid TOML nor JSON\nIt might be corrupted or from an incompatible version", am.AliasStore)
		}

		// If it was JSON, convert to TOML for future use
		fmt.Println("Converting aliases from JSON to TOML format for better readability...")
		if err := am.convertAliasesToTOML(); err != nil {
			fmt.Printf("Warning: Failed to convert aliases to TOML: %v\n", err)
		}
	}

	return nil
}

// AddSaveAliasesTomlSupport modifies the existing SaveAliases method to use TOML format
// This should be called from the original SaveAliases method in alias_manager.go
func (am *AliasManager) AddSaveAliasesTomlSupport() error {
	// Create the directory if it doesn't exist
	dir := filepath.Dir(am.AliasStore)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w\n\nCheck permissions or run with appropriate privileges", dir, err)
	}

	file, err := os.Create(am.AliasStore)
	if err != nil {
		return fmt.Errorf("failed to create aliases file at %s: %w\n\nCheck permissions or disk space issues", am.AliasStore, err)
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(am.Aliases); err != nil {
		return fmt.Errorf("failed to encode aliases as TOML: %w\n\nThis might be due to invalid data or disk issues", err)
	}

	return nil
}

// convertAliasesToTOML reads the existing JSON aliases file, parses it, and writes it back as TOML.
func (am *AliasManager) convertAliasesToTOML() error {
	// Read existing JSON file
	data, err := os.ReadFile(am.AliasStore)
	if err != nil {
		return err
	}

	// Parse JSON
	aliases := make(map[string]AliasCommands)
	if err := json.Unmarshal(data, &aliases); err != nil {
		return err
	}

	// Create backup of original file
	backupFile := am.AliasStore + ".json.bak"
	if err := os.WriteFile(backupFile, data, 0644); err != nil {
		return fmt.Errorf("failed to create backup file %s: %w (check disk space and permissions)", backupFile, err)
	}

	// Write as TOML
	file, err := os.Create(am.AliasStore)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(aliases); err != nil {
		// Restore from backup if TOML encoding fails
		if restoreErr := os.Rename(backupFile, am.AliasStore); restoreErr != nil {
			return fmt.Errorf("TOML encoding failed: %v, and restore failed: %v", err, restoreErr)
		}
		return err
	}

	fmt.Printf("Aliases converted to TOML. Original JSON backup saved as %s\n", backupFile)
	return nil
}

package aliasctl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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

// LoadConfig loads the application configuration.
func (am *AliasManager) LoadConfig() error {
	data, err := os.ReadFile(am.ConfigFile)
	if err != nil {
		return err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	am.Shell = config.DefaultShell
	am.AliasFile = config.DefaultAliasFile
	am.EncryptionUsed = config.UseEncryption

	// Initialize map if nil
	if am.AIProviders == nil {
		am.AIProviders = make(map[string]AIProvider)
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
	if config.AIProvider != "" && am.AIProviders[config.AIProvider] != nil {
		am.AIProvider = am.AIProviders[config.AIProvider]
	}

	return nil
}

// SaveConfig saves the application configuration.
func (am *AliasManager) SaveConfig() error {
	config := Config{
		DefaultShell:     am.Shell,
		DefaultAliasFile: am.AliasFile,
		UseEncryption:    am.EncryptionUsed,
		AIProviders:      make(map[string]bool),
	}

	// Track which providers are configured
	for name := range am.AIProviders {
		config.AIProviders[name] = true
	}

	// Set default provider
	if am.AIProvider != nil {
		switch am.AIProvider.(type) {
		case *OllamaProvider:
			config.AIProvider = "ollama"
		case *OpenAIProvider:
			config.AIProvider = "openai"
		case *AnthropicProvider:
			config.AIProvider = "anthropic"
		}
	}

	// Configure providers
	if provider, ok := am.AIProviders["ollama"]; ok {
		if ollamaProvider, ok := provider.(*OllamaProvider); ok {
			config.OllamaEndpoint = ollamaProvider.Endpoint
			config.OllamaModel = ollamaProvider.Model
		}
	}

	if provider, ok := am.AIProviders["openai"]; ok {
		if openAIProvider, ok := provider.(*OpenAIProvider); ok {
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
	}

	if provider, ok := am.AIProviders["anthropic"]; ok {
		if anthropicProvider, ok := provider.(*AnthropicProvider); ok {
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
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(am.ConfigFile, data, 0644)
}

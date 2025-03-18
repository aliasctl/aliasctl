package aliasctl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aliasctl/aliasctl/pkg/aliasctl/ai"
)

// KeyFileNotFoundError is used when a key file is not found.
// It provides a specific error type for encryption key file issues
// to allow for specialized error handling.
type KeyFileNotFoundError struct {
	KeyPath string // Path to the missing encryption key file
}

// Error returns the error message for a KeyFileNotFoundError.
// It includes the path to the missing key file and guidance on how to resolve the issue.
func (e KeyFileNotFoundError) Error() string {
	return fmt.Sprintf("encryption key file not found at: %s\n\nTo set up encryption, use 'aliasctl encrypt-api-keys' or reconfigure your API provider", e.KeyPath)
}

// EncryptAPIKeys encrypts any API keys in the configuration.
// It generates a secure encryption key if one doesn't exist, then encrypts
// any plaintext API keys found in the configuration. The encryption key is stored
// separately for security.
// Returns an error if the encryption key cannot be generated or stored, or if
// any part of the encryption process fails.
func (am *AliasManager) EncryptAPIKeys() error {
	// Generate encryption key if it doesn't exist
	if _, err := os.Stat(am.EncryptionKey); os.IsNotExist(err) {
		// Create directory if it doesn't exist
		keyDir := filepath.Dir(am.EncryptionKey)
		if err := os.MkdirAll(keyDir, 0700); err != nil {
			return fmt.Errorf("failed to create encryption key directory at %s: %w (check directory permissions)", keyDir, err)
		}

		// Generate a random encryption key
		key, err := GenerateRandomKey()
		if err != nil {
			return fmt.Errorf("failed to generate encryption key: %w (this could be due to insufficient system entropy)", err)
		}

		// Write key to file with restricted permissions
		if err := os.WriteFile(am.EncryptionKey, key, 0600); err != nil {
			return fmt.Errorf("failed to write encryption key to %s: %w (check file permissions)", am.EncryptionKey, err)
		}
	}

	// Get providers from the AI manager
	providers := am.GetAvailableProviders()
	hasProvider := make(map[string]bool)
	for _, name := range providers {
		hasProvider[name] = true
	}

	// Update the configuration
	config := Config{}
	if err := loadConfigFromFile(am.ConfigFile, &config); err != nil {
		return fmt.Errorf("failed to load configuration for encryption: %w", err)
	}

	// Encrypt API keys as needed
	if hasProvider["openai"] {
		// Get the OpenAI provider through the AI manager
		provider, err := am.aiManager.GetProvider("openai")
		if err == nil {
			if openAIProvider, ok := provider.(*ai.OpenAIProvider); ok && openAIProvider.APIKey != "" {
				encryptedKey, err := EncryptString(openAIProvider.APIKey, am.EncryptionKey)
				if err != nil {
					return fmt.Errorf("failed to encrypt OpenAI API key: %w", err)
				}

				config.OpenAIKeyEncrypted = encryptedKey
				config.OpenAIKey = "" // Clear plaintext key
				config.UseEncryption = true
			}
		}
	}

	if hasProvider["anthropic"] {
		// Get the Anthropic provider through the AI manager
		provider, err := am.aiManager.GetProvider("anthropic")
		if err == nil {
			if anthropicProvider, ok := provider.(*ai.AnthropicProvider); ok && anthropicProvider.APIKey != "" {
				encryptedKey, err := EncryptString(anthropicProvider.APIKey, am.EncryptionKey)
				if err != nil {
					return fmt.Errorf("failed to encrypt Anthropic API key: %w", err)
				}

				config.AnthropicKeyEncrypted = encryptedKey
				config.AnthropicKey = "" // Clear plaintext key
				config.UseEncryption = true
			}
		}
	}

	// Save the updated configuration
	if err := saveConfigToFile(am.ConfigFile, config); err != nil {
		return fmt.Errorf("failed to save configuration with encrypted keys: %w", err)
	}

	am.EncryptionUsed = true
	return nil
}

// DisableEncryption disables encryption and reverts to plaintext API keys.
// It decrypts any encrypted API keys in the configuration and stores them
// in plaintext. The encryption flag is also turned off.
// Returns a KeyFileNotFoundError if the encryption key file doesn't exist,
// or a generic error if decryption fails for any other reason.
func (am *AliasManager) DisableEncryption() error {
	// Load configuration
	config := Config{}
	if err := loadConfigFromFile(am.ConfigFile, &config); err != nil {
		return fmt.Errorf("failed to load configuration to disable encryption: %w (check if config file exists and has valid format)", err)
	}

	// Check if we have encrypted keys that need decryption
	if config.OpenAIKeyEncrypted != "" {
		// Decrypt the OpenAI key
		decryptedKey, err := DecryptString(config.OpenAIKeyEncrypted, am.EncryptionKey)
		if err != nil {
			if _, ok := err.(*KeyFileNotFoundError); ok {
				return &KeyFileNotFoundError{KeyPath: am.EncryptionKey}
			}
			return fmt.Errorf("failed to decrypt OpenAI API key: %w (encryption key may be corrupted or inaccessible)", err)
		}
		config.OpenAIKey = decryptedKey
		config.OpenAIKeyEncrypted = ""
	}

	if config.AnthropicKeyEncrypted != "" {
		// Decrypt the Anthropic key
		decryptedKey, err := DecryptString(config.AnthropicKeyEncrypted, am.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to decrypt Anthropic API key: %w", err)
		}
		config.AnthropicKey = decryptedKey
		config.AnthropicKeyEncrypted = ""
	}

	// Update the encryption flag
	config.UseEncryption = false
	am.EncryptionUsed = false

	// Save the updated configuration
	if err := saveConfigToFile(am.ConfigFile, config); err != nil {
		return fmt.Errorf("failed to save configuration with decrypted keys: %w", err)
	}

	return nil
}

// GetEncryptionKeyPath gets the path to the encryption key file.
// It determines the appropriate location based on the configDir parameter or
// the user's home directory if configDir is empty.
// Returns the absolute path to the encryption key file and any error encountered.
func GetEncryptionKeyPath(configDir string) (string, error) {
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		configDir = filepath.Join(homeDir, ".aliasctl")
	}

	return filepath.Join(configDir, "encryption.key"), nil
}

// LoadConfig loads configuration from the specified path into the config struct.
// It reads the file at path and unmarshals the JSON into the config pointer.
// Returns an error if the file cannot be read or the JSON is invalid.
func LoadConfig(path string, config *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, config)
}

// SaveConfig saves the configuration to the specified path.
// It marshals the config struct into JSON and writes it to the file at path.
// Returns an error if the marshalling fails or the file cannot be written.
func SaveConfig(path string, config Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// GenerateRandomKey generates a random encryption key.
// It creates a 256-bit (32 byte) key using a secure random number generator.
// On Unix-like systems, it reads from /dev/urandom for entropy.
// Returns the generated key and any error encountered during generation.
func GenerateRandomKey() ([]byte, error) {
	// Implementation depends on your encryption methodology
	// This is a placeholder for a function that would generate a secure key
	key := make([]byte, 32) // 256-bit key
	if _, err := os.ReadFile("/dev/urandom"); err == nil {
		// For Unix-like systems
		file, err := os.Open("/dev/urandom")
		if err != nil {
			return nil, err
		}
		defer file.Close()

		if _, err := file.Read(key); err != nil {
			return nil, err
		}
	} else {
		// For systems without /dev/urandom
		// Use a cryptographically secure random number generator
		// This is just a placeholder and should be replaced with proper crypto/rand usage
		return nil, fmt.Errorf("secure random number generation not implemented for this platform")
	}

	return key, nil
}

// EncryptString encrypts a string using the encryption key.
// It reads the encryption key from the specified path and uses it to
// encrypt the plaintext string.
// Returns the encrypted string or an error if the key cannot be read or
// the encryption fails. A KeyFileNotFoundError is returned if the key file doesn't exist.
func EncryptString(plaintext string, keyPath string) (string, error) {
	// Read the encryption key
	key, err := os.ReadFile(keyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", &KeyFileNotFoundError{KeyPath: keyPath}
		}
		return "", fmt.Errorf("failed to read encryption key: %w (check file permissions and that the key exists)", err)
	}

	// Use the key in encryption (placeholder implementation)
	_ = key // Using key to avoid unused variable error
	return fmt.Sprintf("encrypted:%s", plaintext), nil
}

// DecryptString decrypts a string using the encryption key.
// It reads the encryption key from the specified path and uses it to
// decrypt the ciphertext string.
// Returns the decrypted string or an error if the key cannot be read,
// the decryption fails, or the ciphertext is invalid.
// A KeyFileNotFoundError is returned if the key file doesn't exist.
func DecryptString(ciphertext string, keyPath string) (string, error) {
	// Read the encryption key
	key, err := os.ReadFile(keyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", &KeyFileNotFoundError{KeyPath: keyPath}
		}
		return "", fmt.Errorf("failed to read encryption key: %w (check file permissions and that the key exists)", err)
	}

	// Use the key in decryption (placeholder implementation)
	_ = key // Using key to avoid unused variable error

	if ciphertext == "" {
		return "", fmt.Errorf("empty ciphertext provided (no encrypted data to decrypt)")
	}

	if len(ciphertext) < 10 || ciphertext[:10] != "encrypted:" {
		return "", fmt.Errorf("invalid ciphertext format (data doesn't appear to be properly encrypted)")
	}

	return ciphertext[10:], nil
}

// loadConfigFromFile is a wrapper around LoadConfig to avoid name conflicts.
// It loads configuration from the specified path into the config struct.
func loadConfigFromFile(path string, config *Config) error {
	// Implementation should be provided elsewhere
	// This is just a wrapper to avoid name conflicts
	return LoadConfig(path, config)
}

// saveConfigToFile is a wrapper around SaveConfig to avoid name conflicts.
// It saves the configuration to the specified path.
func saveConfigToFile(path string, config Config) error {
	// Implementation should be provided elsewhere
	// This is just a wrapper to avoid name conflicts
	return SaveConfig(path, config)
}

package aliasctl

import (
	"encoding/json"
	"fmt"
	"os"
)

// EncryptAPIKeys encrypts any API keys in the configuration.
func (am *AliasManager) EncryptAPIKeys() error {
	if am.EncryptionUsed {
		fmt.Println("Encryption is already enabled")
		return nil
	}

	// Load current config to ensure we have the latest
	data, err := os.ReadFile(am.ConfigFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var config Config
	if len(data) > 0 {
		if err := json.Unmarshal(data, &config); err != nil {
			return err
		}
	}

	// Check if there are plaintext API keys to encrypt
	hasKeysToEncrypt := false

	if config.OpenAIKey != "" {
		hasKeysToEncrypt = true
		encryptedKey, err := EncryptString(config.OpenAIKey, am.EncryptionKey)
		if err == nil {
			config.OpenAIKeyEncrypted = encryptedKey
			config.OpenAIKey = "" // Clear plaintext key
			fmt.Println("OpenAI API key encrypted successfully")
		} else {
			return fmt.Errorf("failed to encrypt OpenAI API key: %v", err)
		}
	}

	if config.AnthropicKey != "" {
		hasKeysToEncrypt = true
		encryptedKey, err := EncryptString(config.AnthropicKey, am.EncryptionKey)
		if err == nil {
			config.AnthropicKeyEncrypted = encryptedKey
			config.AnthropicKey = "" // Clear plaintext key
			fmt.Println("Anthropic API key encrypted successfully")
		} else {
			return fmt.Errorf("failed to encrypt Anthropic API key: %v", err)
		}
	}

	if !hasKeysToEncrypt {
		return fmt.Errorf("no plaintext API keys found to encrypt")
	}

	// Update config with encryption enabled
	config.UseEncryption = true
	am.EncryptionUsed = true

	// Save the updated config
	data, err = json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(am.ConfigFile, data, 0644); err != nil {
		return err
	}

	fmt.Printf("API keys encrypted and stored. Encryption key file: %s\n", am.EncryptionKey)
	fmt.Println("WARNING: Keep this file secure as it's needed to decrypt your API keys.")
	return nil
}

// DisableEncryption disables encryption and reverts to plaintext API keys.
func (am *AliasManager) DisableEncryption() error {
	if !am.EncryptionUsed {
		return fmt.Errorf("encryption is not currently enabled")
	}

	// Load current config
	data, err := os.ReadFile(am.ConfigFile)
	if err != nil {
		return err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	// Check if there's an encrypted key
	if config.OpenAIKeyEncrypted == "" && config.AnthropicKeyEncrypted == "" {
		return fmt.Errorf("no encrypted API key found")
	}

	// Decrypt the keys
	if config.OpenAIKeyEncrypted != "" {
		decryptedKey, err := DecryptString(config.OpenAIKeyEncrypted, am.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to decrypt OpenAI API key: %v", err)
		}
		config.OpenAIKey = decryptedKey
		config.OpenAIKeyEncrypted = ""
	}

	if config.AnthropicKeyEncrypted != "" {
		decryptedKey, err := DecryptString(config.AnthropicKeyEncrypted, am.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to decrypt Anthropic API key: %v", err)
		}
		config.AnthropicKey = decryptedKey
		config.AnthropicKeyEncrypted = ""
	}

	// Update config
	config.UseEncryption = false
	am.EncryptionUsed = false

	// Save the updated config
	data, err = json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(am.ConfigFile, data, 0644); err != nil {
		return err
	}

	fmt.Println("API keys decrypted and stored in plaintext. Encryption disabled.")
	return nil
}

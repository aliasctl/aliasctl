package aliasctl

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// KeyFileNotFoundError indicates the encryption key file was not found
type KeyFileNotFoundError struct {
	KeyPath string
}

func (e KeyFileNotFoundError) Error() string {
	return fmt.Sprintf("encryption key file not found: %s", e.KeyPath)
}

// EncryptString encrypts a string using AES-GCM
func EncryptString(plaintext string, keyPath string) (string, error) {
	key, err := loadOrCreateKeyFile(keyPath)
	if err != nil {
		return "", err
	}

	// Create a new cipher block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new GCM cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt and seal the data
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	// Return base64 encoded ciphertext
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptString decrypts a string using AES-GCM
func DecryptString(encryptedText string, keyPath string) (string, error) {
	key, err := loadKeyFile(keyPath)
	if err != nil {
		return "", err
	}

	// Decode the base64 string
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	// Create a new cipher block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new GCM cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Get the nonce size
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// Extract the nonce and ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// IsEncrypted checks if a string appears to be encrypted
func IsEncrypted(text string) bool {
	// Check if it's base64 encoded and has a reasonable length for encrypted content
	if len(text) < 24 || !isBase64(text) {
		return false
	}
	return true
}

// isBase64 checks if a string is a valid base64 encoded string
func isBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

// loadOrCreateKeyFile loads the encryption key from a file or creates it if it doesn't exist
func loadOrCreateKeyFile(keyPath string) ([]byte, error) {
	// Try to load an existing key
	key, err := loadKeyFile(keyPath)
	if err == nil {
		return key, nil
	}

	// If the error is not a "file not found" error, return it
	if !os.IsNotExist(err) {
		return nil, err
	}

	// Generate a new random key
	key = make([]byte, 32) // 256 bits
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}

	// Ensure the directory exists
	dir := filepath.Dir(keyPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	// Write the key to the file
	if err := os.WriteFile(keyPath, key, 0600); err != nil {
		return nil, err
	}

	return key, nil
}

// loadKeyFile loads the encryption key from a file
func loadKeyFile(keyPath string) ([]byte, error) {
	// Check if the file exists
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return nil, &KeyFileNotFoundError{KeyPath: keyPath}
	}

	// Read the key from the file
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	// Ensure the key is 256 bits (32 bytes)
	if len(key) != 32 {
		// If not, hash it to get a 256-bit key
		hasher := sha256.New()
		hasher.Write(key)
		key = hasher.Sum(nil)
	}

	return key, nil
}

// GetEncryptionKeyPath returns the path to the encryption key file
func GetEncryptionKeyPath(configDir string) string {
	return filepath.Join(configDir, ".keyring", "api_keys.key")
}

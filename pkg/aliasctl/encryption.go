package aliasctl

import (
	"fmt"
)

// LegacyEncryptString provides backward compatibility for encryption
func LegacyEncryptString(plaintext string) (string, error) {
	// Implementation
	return fmt.Sprintf("encrypted:%s", plaintext), nil
}

// LegacyDecryptString provides backward compatibility for decryption
func LegacyDecryptString(ciphertext string) (string, error) {
	// Implementation
	if ciphertext == "" {
		return "", fmt.Errorf("empty ciphertext provided")
	}

	if len(ciphertext) < 10 || ciphertext[:10] != "encrypted:" {
		return "", fmt.Errorf("invalid ciphertext format")
	}

	return ciphertext[10:], nil
}

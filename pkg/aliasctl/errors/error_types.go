package errors

import (
	"fmt"
)

// StandardizedError represents a standardized error with hints.
// It includes a message, optional cause error, and suggestions for resolving the issue.
type StandardizedError struct {
	Message string   // The main error message
	Cause   error    // The underlying error that caused this one, if any
	Hints   []string // Suggestions for resolving the error
}

// Error returns the error message with hints.
// It formats a complete error message including the cause and any hints/suggestions.
func (e *StandardizedError) Error() string {
	result := e.Message
	if e.Cause != nil {
		result += ": " + e.Cause.Error()
	}

	if len(e.Hints) > 0 {
		result += "\n\n"
		for _, hint := range e.Hints {
			result += "- " + hint + "\n"
		}
	}

	return result
}

// Unwrap returns the underlying error.
// This allows errors.Is() and errors.As() to work with wrapped errors.
func (e *StandardizedError) Unwrap() error {
	return e.Cause
}

// NetworkError represents a network connectivity error.
// It includes the endpoint that couldn't be reached and the cause of the failure.
type NetworkError struct {
	Endpoint string // The URL or address that couldn't be connected to
	Cause    error  // The underlying error that caused the network failure
}

// Error returns the error message.
// It formats a complete error message with suggestions for resolving network issues.
func (e *NetworkError) Error() string {
	msg := fmt.Sprintf("failed to connect to %s", e.Endpoint)
	if e.Cause != nil {
		msg += ": " + e.Cause.Error()
	}

	msg += "\n\nPossible solutions:"
	msg += "\n- Check that the service is running"
	msg += "\n- Verify network connectivity"
	msg += "\n- Confirm the endpoint URL is correct"

	if e.Endpoint == "http://localhost:11434" {
		msg += "\n- If using Ollama, ensure it's started with 'ollama serve'"
	}

	return msg
}

// Unwrap returns the underlying error.
// This allows errors.Is() and errors.As() to work with wrapped errors.
func (e *NetworkError) Unwrap() error {
	return e.Cause
}

// PermissionError represents a file permission error.
// It includes the path to the file/directory and the cause of the permission issue.
type PermissionError struct {
	Path  string // The path to the file or directory with permission issues
	Cause error  // The underlying error that caused the permission failure
}

// Error returns the error message.
// It formats a complete error message with suggestions for resolving permission issues.
func (e *PermissionError) Error() string {
	msg := fmt.Sprintf("permission denied for %s", e.Path)
	if e.Cause != nil {
		msg += ": " + e.Cause.Error()
	}

	msg += "\n\nPossible solutions:"
	msg += "\n- Check file/directory permissions"
	msg += "\n- Verify the directory exists"
	msg += "\n- Try specifying a different location with 'aliasctl set-file'"

	return msg
}

// Unwrap returns the underlying error.
// This allows errors.Is() and errors.As() to work with wrapped errors.
func (e *PermissionError) Unwrap() error {
	return e.Cause
}

// ConfigurationError represents an error in configuration.
// It includes the component with the configuration issue, the cause, and optional hints.
type ConfigurationError struct {
	Component string   // The component or section with the configuration issue
	Cause     error    // The underlying error that caused the configuration failure
	Hints     []string // Custom hints for resolving this specific configuration error
}

// Error returns the error message.
// It formats a complete error message with suggestions for resolving configuration issues.
func (e *ConfigurationError) Error() string {
	msg := fmt.Sprintf("configuration error in %s", e.Component)
	if e.Cause != nil {
		msg += ": " + e.Cause.Error()
	}

	msg += "\n\nPossible solutions:"
	if len(e.Hints) > 0 {
		for _, hint := range e.Hints {
			msg += "\n- " + hint
		}
	} else {
		msg += "\n- Check configuration file format"
		msg += "\n- Ensure configuration directory exists"
		msg += "\n- Try reconfiguring this component"
	}

	return msg
}

// Unwrap returns the underlying error.
// This allows errors.Is() and errors.As() to work with wrapped errors.
func (e *ConfigurationError) Unwrap() error {
	return e.Cause
}

// NotFoundError represents an error when a resource is not found.
// It includes the type and name of the resource that wasn't found, along with an optional hint.
type NotFoundError struct {
	ResourceType string // The type of resource (e.g., "alias", "file", "directory")
	ResourceName string // The name of the resource that wasn't found
	Hint         string // A specific suggestion for resolving the issue
}

// Error returns the error message.
// It formats a complete error message with the hint if provided.
func (e *NotFoundError) Error() string {
	msg := fmt.Sprintf("%s '%s' not found", e.ResourceType, e.ResourceName)

	if e.Hint != "" {
		msg += "\n\n" + e.Hint
	}

	return msg
}

// KeyFileNotFoundError represents when an encryption key file is missing.
// It includes the path to the key file that wasn't found.
type KeyFileNotFoundError struct {
	Path string // The path to the encryption key file that wasn't found
}

// Error returns the error message.
// It formats a complete error message with a suggestion for resolving the issue.
func (e *KeyFileNotFoundError) Error() string {
	return fmt.Sprintf("encryption key file not found at: %s\n\nTo set up encryption, use 'aliasctl encrypt-api-keys' or reconfigure your API provider", e.Path)
}

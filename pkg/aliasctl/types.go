package aliasctl

// ShellType represents the type of shell.
type ShellType string

const (
	// Shell types
	ShellBash           ShellType = "bash"
	ShellZsh            ShellType = "zsh"
	ShellFish           ShellType = "fish"
	ShellKsh            ShellType = "ksh"
	ShellPowerShell     ShellType = "powershell"
	ShellPowerShellCore ShellType = "pwsh"
	ShellCmd            ShellType = "cmd"
)

// AliasCommands holds the commands for all supported shells.
type AliasCommands struct {
	Bash           string `json:"bash"`
	Zsh            string `json:"zsh"`
	Fish           string `json:"fish"`
	Ksh            string `json:"ksh"`
	PowerShell     string `json:"powershell"`
	PowerShellCore string `json:"pwsh"`
	Cmd            string `json:"cmd"`
}

// AliasManager handles platform-specific alias operations.
type AliasManager struct {
	Platform       string                   // The operating system platform
	Shell          ShellType                // The type of shell
	AliasFile      string                   // The path to the alias file
	Aliases        map[string]AliasCommands // A map of alias names to shell-specific commands
	AIProvider     AIProvider               // The configured AI provider (for backward compatibility)
	AIProviders    map[string]AIProvider    // Map of configured AI providers by name
	AIConfigured   bool                     // Whether an AI provider is configured
	ConfigDir      string                   // The configuration directory
	AliasStore     string                   // The path to the alias store file
	ConfigFile     string                   // The path to the configuration file
	EncryptionKey  string                   // The path to the encryption key file
	EncryptionUsed bool                     // Whether encryption is being used
}

// Config represents the application configuration.
type Config struct {
	DefaultShell          ShellType       `json:"default_shell"`           // The default shell type
	DefaultAliasFile      string          `json:"default_alias_file"`      // The default alias file path
	AIProvider            string          `json:"ai_provider"`             // The default AI provider type
	AIProviders           map[string]bool `json:"ai_providers"`            // Map of configured AI providers
	OllamaEndpoint        string          `json:"ollama_endpoint"`         // The Ollama endpoint URL
	OllamaModel           string          `json:"ollama_model"`            // The Ollama model name
	OpenAIEndpoint        string          `json:"openai_endpoint"`         // The OpenAI endpoint URL
	OpenAIKey             string          `json:"openai_key"`              // The OpenAI API key (plaintext, deprecated)
	OpenAIKeyEncrypted    string          `json:"openai_key_encrypted"`    // The OpenAI API key (encrypted)
	OpenAIModel           string          `json:"openai_model"`            // The OpenAI model name
	AnthropicEndpoint     string          `json:"anthropic_endpoint"`      // The Anthropic endpoint URL
	AnthropicKey          string          `json:"anthropic_key"`           // The Anthropic API key (plaintext, deprecated)
	AnthropicKeyEncrypted string          `json:"anthropic_key_encrypted"` // The Anthropic API key (encrypted)
	AnthropicModel        string          `json:"anthropic_model"`         // The Anthropic model name
	UseEncryption         bool            `json:"use_encryption"`          // Whether to use encryption for API keys
}

// AIProvider interface for AI services.
type AIProvider interface {
	ConvertAlias(alias, fromShell, toShell string) (string, error) // Converts an alias from one shell to another
}

// OllamaProvider implements AIProvider for Ollama.
type OllamaProvider struct {
	Endpoint string // The Ollama endpoint URL
	Model    string // The Ollama model name
}

// OpenAIProvider implements AIProvider for OpenAI-compatible APIs.
type OpenAIProvider struct {
	Endpoint string // The OpenAI endpoint URL
	APIKey   string // The OpenAI API key
	Model    string // The OpenAI model name
}

// AnthropicProvider implements AIProvider for Anthropic Claude.
type AnthropicProvider struct {
	Endpoint string // The Anthropic endpoint URL
	APIKey   string // The Anthropic API key
	Model    string // The Anthropic model name
}

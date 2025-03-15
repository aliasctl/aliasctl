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

// AliasManager handles platform-specific alias operations.
type AliasManager struct {
	Platform     string            // The operating system platform
	Shell        ShellType         // The type of shell
	AliasFile    string            // The path to the alias file
	Aliases      map[string]string // A map of alias names to commands
	AIProvider   AIProvider        // The configured AI provider
	AIConfigured bool              // Whether an AI provider is configured
	ConfigDir    string            // The configuration directory
	AliasStore   string            // The path to the alias store file
	ConfigFile   string            // The path to the configuration file
}

// Config represents the application configuration.
type Config struct {
	DefaultShell     ShellType `json:"default_shell"`      // The default shell type
	DefaultAliasFile string    `json:"default_alias_file"` // The default alias file path
	AIProvider       string    `json:"ai_provider"`        // The AI provider type
	OllamaEndpoint   string    `json:"ollama_endpoint"`    // The Ollama endpoint URL
	OllamaModel      string    `json:"ollama_model"`       // The Ollama model name
	OpenAIEndpoint   string    `json:"openai_endpoint"`    // The OpenAI endpoint URL
	OpenAIKey        string    `json:"openai_key"`         // The OpenAI API key
	OpenAIModel      string    `json:"openai_model"`       // The OpenAI model name
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

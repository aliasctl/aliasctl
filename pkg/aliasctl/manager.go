package aliasctl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// NewAliasManager creates a new AliasManager.
func NewAliasManager() *AliasManager {
	platform := runtime.GOOS
	configDir := getConfigDir()

	am := &AliasManager{
		Platform:     platform,
		Aliases:      make(map[string]string),
		AIConfigured: false,
		ConfigDir:    configDir,
		AliasStore:   filepath.Join(configDir, "aliases.json"),
		ConfigFile:   filepath.Join(configDir, "config.json"),
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

	if config.AIProvider == "ollama" && config.OllamaEndpoint != "" && config.OllamaModel != "" {
		am.ConfigureOllama(config.OllamaEndpoint, config.OllamaModel)
	} else if config.AIProvider == "openai" && config.OpenAIEndpoint != "" && config.OpenAIKey != "" && config.OpenAIModel != "" {
		am.ConfigureOpenAI(config.OpenAIEndpoint, config.OpenAIKey, config.OpenAIModel)
	}

	return nil
}

// SaveConfig saves the application configuration.
func (am *AliasManager) SaveConfig() error {
	config := Config{
		DefaultShell:     am.Shell,
		DefaultAliasFile: am.AliasFile,
	}

	if am.AIConfigured {
		switch provider := am.AIProvider.(type) {
		case *OllamaProvider:
			config.AIProvider = "ollama"
			config.OllamaEndpoint = provider.Endpoint
			config.OllamaModel = provider.Model
		case *OpenAIProvider:
			config.AIProvider = "openai"
			config.OpenAIEndpoint = provider.Endpoint
			config.OpenAIKey = provider.APIKey
			config.OpenAIModel = provider.Model
		}
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(am.ConfigFile, data, 0644)
}

// LoadAliases loads aliases from the JSON store.
func (am *AliasManager) LoadAliases() error {
	if _, err := os.Stat(am.AliasStore); os.IsNotExist(err) {
		am.Aliases = make(map[string]string)
		return nil
	}

	data, err := os.ReadFile(am.AliasStore)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &am.Aliases)
}

// SaveAliases saves aliases to the JSON store.
func (am *AliasManager) SaveAliases() error {
	data, err := json.MarshalIndent(am.Aliases, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(am.AliasStore, data, 0644)
}

// AddAlias adds a new alias.
func (am *AliasManager) AddAlias(name, command string) {
	am.Aliases[name] = command
}

// RemoveAlias removes an alias.
func (am *AliasManager) RemoveAlias(name string) bool {
	if _, exists := am.Aliases[name]; exists {
		delete(am.Aliases, name)
		return true
	}
	return false
}

// ListAliases prints all aliases.
func (am *AliasManager) ListAliases() {
	fmt.Printf("Aliases for %s shell on %s platform:\n", am.Shell, am.Platform)
	if len(am.Aliases) == 0 {
		fmt.Println("No aliases defined.")
		return
	}

	for name, command := range am.Aliases {
		fmt.Printf("%s = %s\n", name, command)
	}
}

// SetShell manually sets the shell type.
func (am *AliasManager) SetShell(shell string) error {
	switch shell {
	case "bash":
		am.Shell = ShellBash
	case "zsh":
		am.Shell = ShellZsh
	case "fish":
		am.Shell = ShellFish
	case "ksh":
		am.Shell = ShellKsh
	case "powershell":
		am.Shell = ShellPowerShell
	case "pwsh":
		am.Shell = ShellPowerShellCore
	case "cmd":
		am.Shell = ShellCmd
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}
	return am.SaveConfig()
}

// SetAliasFile manually sets the alias file path.
func (am *AliasManager) SetAliasFile(filePath string) error {
	am.AliasFile = filePath
	return am.SaveConfig()
}

// ConfigureOllama sets up the Ollama AI provider.
func (am *AliasManager) ConfigureOllama(endpoint, model string) {
	am.AIProvider = &OllamaProvider{
		Endpoint: endpoint,
		Model:    model,
	}
	am.AIConfigured = true
	am.SaveConfig()
}

// ConfigureOpenAI sets up the OpenAI-compatible AI provider.
func (am *AliasManager) ConfigureOpenAI(endpoint, apiKey, model string) {
	am.AIProvider = &OpenAIProvider{
		Endpoint: endpoint,
		APIKey:   apiKey,
		Model:    model,
	}
	am.AIConfigured = true
	am.SaveConfig()
}

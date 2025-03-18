package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all configuration values for the application
// including paths to relevant files and directories, as well as user preferences.
type Config struct {
	AliasFile     string `mapstructure:"alias_file"`     // Path to the shell's main alias file
	AliasesDir    string `mapstructure:"aliases_dir"`    // Directory where individual alias files are stored
	ConfigFile    string `mapstructure:"config_file"`    // Path to aliasctl's own configuration file
	DefaultEditor string `mapstructure:"default_editor"` // Default editor to use when not specified by EDITOR env var
}

// LoadConfig loads the configuration using Viper.
// It sets default values, reads from the config file if available,
// and creates necessary directories. If the config file doesn't exist,
// defaults will be used.
func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("unable to determine user home directory: %w", err)
	}

	// Set default config values
	viper.SetDefault("alias_file", filepath.Join(home, ".zshrc"))
	viper.SetDefault("aliases_dir", filepath.Join(home, ".aliases"))
	viper.SetDefault("config_file", filepath.Join(home, ".aliasctl.yaml"))
	viper.SetDefault("default_editor", "vim")

	// Set config name, paths, and type
	viper.SetConfigName(".aliasctl")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(home)
	viper.AddConfigPath(".")

	// Read config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Only return error if it's not just missing config file
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
		// Config file not found - we'll use defaults
	}

	// Parse config into struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Create aliases directory if it doesn't exist
	if err := os.MkdirAll(config.AliasesDir, 0755); err != nil {
		return nil, fmt.Errorf("unable to create aliases directory: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the current configuration to file.
// It updates all configuration values in Viper and writes
// them to the config file specified in the Config struct.
func SaveConfig(config *Config) error {
	viper.Set("alias_file", config.AliasFile)
	viper.Set("aliases_dir", config.AliasesDir)
	viper.Set("config_file", config.ConfigFile)
	viper.Set("default_editor", config.DefaultEditor)

	return viper.WriteConfig()
}

// InitConfig creates a new config file if it doesn't exist.
// It loads default configuration values and saves them to
// create an initial configuration file with sensible defaults.
// If the config file already exists, it does nothing.
func InitConfig() error {
	if _, err := os.Stat(viper.ConfigFileUsed()); os.IsNotExist(err) {
		config, err := LoadConfig()
		if err != nil {
			return err
		}

		return SaveConfig(config)
	}
	return nil
}

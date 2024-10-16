package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	TemplatesDir string `yaml:"templates_dir"`
}

var AppConfig Config

func LoadConfig() error {
	// Set default values
	AppConfig = Config{
		TemplatesDir: "templates",
	}

	configPath := getConfigPath()
	if configPath == "" {
		// Config file doesn't exist, create a default one
		err := createDefaultConfig()
		if err != nil {
			return fmt.Errorf("failed to create default config: %w", err)
		}
		configPath = ".md_config.yaml"
	}

	// Read from config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(data, &AppConfig)
	if err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Override with environment variables
	if envTemplatesDir := os.Getenv("MD_TEMPLATES_DIR"); envTemplatesDir != "" {
		AppConfig.TemplatesDir = envTemplatesDir
	}

	return nil
}

func getConfigPath() string {
	configFileName := ".md_config.yaml"

	// Check for config file in the current directory
	if _, err := os.Stat(configFileName); err == nil {
		return configFileName
	}

	// Check for config file in the user's home directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(homeDir, configFileName)
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
	}

	return ""
}

func createDefaultConfig() error {
	configContent := `# Markdown File Creator Configuration

# templates_dir: Directory where template files are stored
# Default is 'templates' in the current working directory
# Examples for custom paths:
#   Windows: C:\Users\YourUsername\Documents\templates
#   macOS:   /Users/YourUsername/Documents/templates
#   Linux:   /home/YourUsername/Documents/templates

templates_dir: templates

# Note: You can use the environment variable MD_TEMPLATES_DIR to override
# the templates directory at runtime.
`

	err := os.WriteFile(".md_config.yaml", []byte(configContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write default config file: %w", err)
	}

	fmt.Println("Created default configuration file: .md_config.yaml")
	return nil
}

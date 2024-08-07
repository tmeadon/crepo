package main

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	RepoRoot   string   `yaml:"repo_root"`
	Exclusions []string `yaml:"exclusions"`
}

func getDefaultConfig() *Config {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatalf("Error getting user home directory: %v", err)
	}

	return &Config{
		RepoRoot: homeDir,
		Exclusions: []string{
			"node_modules",
			"vendor",
		},
	}
}

func getConfigFilePath() (string, error) {
	configPath := os.Getenv("CREPO_CONFIG_PATH")
	if configPath != "" {
		return configPath, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".crepo")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.yaml"), nil
}

func loadConfig() (*Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		conf := getDefaultConfig()

		if err := saveConfig(conf); err != nil {
			log.Fatalf("Error saving default config: %v", err)
		}

		return conf, nil
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func saveConfig(config *Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	if err := os.WriteFile(configFilePath, data, 0644); err != nil {
		return err
	}

	return nil
}

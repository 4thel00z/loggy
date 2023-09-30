package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	ErrorDescriptions map[string]string `json:"error_descriptions"`
	DatabasePath      string            `json:"database_path"`
}

type LogEntry struct {
	Key         string `json:"key"`
	Message     string `json:"message"`
	Environment string `json:"environment"`
	AppVersion  string `json:"app_version"`
}

type ProblemDetails struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", nil
	}
	configDir := filepath.Join(homeDir, ".config", "loggy")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", nil
	}
	return filepath.Join(configDir, "config.json"), nil

}
func ensureConfig(configFile string) (*Config, error) {

	config := &Config{
		ErrorDescriptions: map[string]string{
			"bad-request":           "### Bad Request\nThe server could not understand the request due to invalid syntax.",
			"internal-server-error": "### Internal Server Error\nThe server encountered an internal error and was unable to complete your request.",
		},
		DatabasePath: filepath.Join(os.TempDir(), "loggy", "logs.db"),
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		configData, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(configFile, configData, 0644); err != nil {
			return nil, err
		}
		return config, nil
	}

	configData, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(configData, config); err != nil {
		return nil, err
	}

	return config, nil
}

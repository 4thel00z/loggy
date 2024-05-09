package loggy

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	ErrorDescriptions map[string]string `json:"error_descriptions"`
	DatabasePath      string            `json:"database_path"`
}

type LogEntry struct {
	Key         string    `json:"key" db:"key"`
	Message     string    `json:"message" db:"message"`
	Environment string    `json:"environment" db:"environment"`
	AppVersion  string    `json:"app_version" db:"app_version"`
	DeviceName  string    `json:"device_name" db:"device_name"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

func (l LogEntry) Title() string {
	return fmt.Sprintf("[%s] %s - %s@%s %sv: %s", l.Key, l.CreatedAt.Format(time.DateTime), l.DeviceName, l.Environment, l.AppVersion)
}
func (l LogEntry) Description() string { return l.Message }

type LogEntries []LogEntry

func (l LogEntry) FilterValue() string {
	return l.String()
}

func (l LogEntries) ToItems() []list.Item {
	var items []list.Item
	for _, i2 := range l {
		items = append(items, i2)
	}
	return items
}

func (l LogEntry) String() string {
	return fmt.Sprintf(
		"[%s] %s %s %s: %s",
		l.Environment, l.AppVersion, l.DeviceName, l.Key, l.Message,
	)
}

type ProblemDetails struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func GetConfigPath() (string, error) {
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
func EnsureConfig(configFile string) (*Config, error) {

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

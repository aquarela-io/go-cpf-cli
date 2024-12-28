package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/posthog/posthog-go"
)

// Config represents telemetry configuration
type Config struct {
	Enabled bool `json:"enabled"`
}

// Event represents a telemetry event
type Event struct {
	Command   string            `json:"command"`
	Success   bool              `json:"success"`
	Error     string           `json:"error,omitempty"`
	OS        string           `json:"os"`
	Arch      string           `json:"arch"`
	Version   string           `json:"version"`
	Timestamp time.Time        `json:"timestamp"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

var (
	config     *Config
	configPath string
	version    string // Will be set during initialization
	apiKey     string // Will be set at build time
	client     posthog.Client
)

// Initialize sets up telemetry with the given version
func Initialize(v string) error {
	version = v

	// Initialize PostHog client if we have an API key
	if apiKey != "" {
		var err error
		client, err = posthog.NewWithConfig(apiKey, posthog.Config{
			Endpoint: "https://us.i.posthog.com",
		})
		if err != nil {
			return fmt.Errorf("failed to initialize PostHog client: %w", err)
		}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".cpf-cli")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath = filepath.Join(configDir, "telemetry.json")
	config = &Config{}

	// Try to load existing config
	if err := loadConfig(); err != nil {
		// If file doesn't exist, create default config (disabled by default)
		config.Enabled = false
		if err := saveConfig(); err != nil {
			return fmt.Errorf("failed to save default config: %w", err)
		}
	}

	return nil
}

// loadConfig loads the telemetry configuration from disk
func loadConfig() error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, config)
}

// saveConfig saves the telemetry configuration to disk
func saveConfig() error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

// SetEnabled enables or disables telemetry
func SetEnabled(enabled bool) error {
	config.Enabled = enabled
	return saveConfig()
}

// IsEnabled returns whether telemetry is enabled
func IsEnabled() bool {
	return config != nil && config.Enabled && apiKey != "" && client != nil
}

// Close closes the PostHog client
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}

// Track sends a telemetry event if telemetry is enabled
func Track(command string, success bool, err error, metadata map[string]string) {
	if !IsEnabled() {
		return
	}

	// Create a unique identifier for the installation
	distinctId := fmt.Sprintf("%s-%s-%s", runtime.GOOS, runtime.GOARCH, version)

	// Convert metadata to interface{} map
	properties := make(map[string]interface{})
	properties["command"] = command
	properties["success"] = success
	properties["os"] = runtime.GOOS
	properties["arch"] = runtime.GOARCH
	properties["version"] = version
	properties["timestamp"] = time.Now().UTC()

	if err != nil {
		properties["error"] = err.Error()
	}

	for k, v := range metadata {
		properties[k] = v
	}

	// Send event asynchronously
	client.Enqueue(posthog.Capture{
		DistinctId: distinctId,
		Event:      "cli_command",
		Properties: properties,
	})
} 
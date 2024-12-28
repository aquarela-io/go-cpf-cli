package telemetry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Config represents telemetry configuration
type Config struct {
	Enabled bool `json:"enabled"`
}

// PostHogEvent represents the event structure for PostHog
type PostHogEvent struct {
	ApiKey      string                 `json:"api_key"`
	Event       string                 `json:"event"`
	DistinctId  string                 `json:"distinct_id"`
	Properties  map[string]interface{} `json:"properties"`
	Timestamp   string                 `json:"timestamp"`
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
	// PostHog Cloud endpoint
	posthogEndpoint = "https://app.posthog.com/capture"
)

// Initialize sets up telemetry with the given version
func Initialize(v string) error {
	version = v
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
	return config != nil && config.Enabled && apiKey != ""
}

// Track sends a telemetry event if telemetry is enabled
func Track(command string, success bool, err error, metadata map[string]string) {
	if !IsEnabled() {
		return
	}

	event := Event{
		Command:   command,
		Success:   success,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		Version:   version,
		Timestamp: time.Now().UTC(),
		Metadata:  metadata,
	}

	if err != nil {
		event.Error = err.Error()
	}

	// Send event asynchronously to not block the main execution
	go func() {
		if err := sendEvent(event); err != nil {
			// Silently fail - we don't want telemetry errors to affect the user
			_ = err
		}
	}()
}

// sendEvent sends the telemetry event to PostHog
func sendEvent(event Event) error {
	// Convert our event to PostHog format
	properties := map[string]interface{}{
		"command":   event.Command,
		"success":   event.Success,
		"os":        event.OS,
		"arch":      event.Arch,
		"version":   event.Version,
		"timestamp": event.Timestamp,
	}

	if event.Error != "" {
		properties["error"] = event.Error
	}

	for k, v := range event.Metadata {
		properties[k] = v
	}

	phEvent := PostHogEvent{
		ApiKey:     apiKey,
		Event:      "cli_command",
		DistinctId: event.OS + "-" + event.Arch, // Anonymous identifier
		Properties: properties,
		Timestamp:  event.Timestamp.Format(time.RFC3339),
	}

	data, err := json.Marshal(phEvent)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", posthogEndpoint, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
} 
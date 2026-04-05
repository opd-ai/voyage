package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ConfigPath returns the platform-appropriate path for configuration files.
// Returns ~/.config/voyage on Linux, ~/Library/Application Support/voyage on macOS,
// and %APPDATA%/voyage on Windows.
func ConfigPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fall back to home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ".voyage"
		}
		return filepath.Join(homeDir, ".voyage")
	}
	return filepath.Join(configDir, "voyage")
}

// ConfigFilePath returns the full path to the config file.
func ConfigFilePath() string {
	return filepath.Join(ConfigPath(), "config.json")
}

// SaveConfig persists the configuration to disk.
func SaveConfig(cfg *Config) error {
	path := ConfigPath()
	if err := os.MkdirAll(path, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ConfigFilePath(), data, 0o644)
}

// LoadConfig loads the configuration from disk.
// Returns default config if file doesn't exist.
func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(ConfigFilePath())
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, err
	}

	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return DefaultConfig(), nil
	}

	// Ensure input config is initialized
	if cfg.Input == nil {
		cfg.Input = NewInputConfig()
	}

	return cfg, nil
}

// ConfigExists returns true if a config file exists.
func ConfigExists() bool {
	_, err := os.Stat(ConfigFilePath())
	return err == nil
}

// DeleteConfig removes the configuration file.
func DeleteConfig() error {
	return os.Remove(ConfigFilePath())
}

// InputBindingData represents serializable input binding data.
type InputBindingData struct {
	Action  int    `json:"action"`
	KeyCode int    `json:"keyCode"`
	KeyName string `json:"keyName"`
}

// MarshalJSON implements custom JSON marshaling for InputConfig.
func (ic *InputConfig) MarshalJSON() ([]byte, error) {
	bindings := make([]InputBindingData, 0)
	for action, binding := range ic.bindings {
		bindings = append(bindings, InputBindingData{
			Action:  int(action),
			KeyCode: binding.KeyCode,
			KeyName: binding.KeyName,
		})
	}
	return json.Marshal(bindings)
}

// UnmarshalJSON implements custom JSON unmarshaling for InputConfig.
func (ic *InputConfig) UnmarshalJSON(data []byte) error {
	var bindings []InputBindingData
	if err := json.Unmarshal(data, &bindings); err != nil {
		return err
	}

	ic.bindings = make(map[Action]KeyBinding)
	for _, b := range bindings {
		ic.bindings[Action(b.Action)] = KeyBinding{
			Action:  Action(b.Action),
			KeyCode: b.KeyCode,
			KeyName: b.KeyName,
		}
	}

	// Fill in any missing defaults
	defaults := NewInputConfig()
	for action, binding := range defaults.bindings {
		if _, exists := ic.bindings[action]; !exists {
			ic.bindings[action] = binding
		}
	}

	return nil
}

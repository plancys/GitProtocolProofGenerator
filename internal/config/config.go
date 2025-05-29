package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the configuration for the report generator
type Config struct {
	// Header template configuration
	Header HeaderConfig `json:"header"`

	// PDF styling configuration
	PDF PDFConfig `json:"pdf"`
}

// HeaderConfig contains the configurable header template
type HeaderConfig struct {
	// Template for the header with placeholders
	Template string `json:"template"`

	// Default executor information
	ExecutorName  string `json:"executor_name"`
	ExecutorEmail string `json:"executor_email"`

	// Default recipient information
	RecipientName string `json:"recipient_name"`

	// Location for the report
	Location string `json:"location"`
}

// PDFConfig contains PDF styling options
type PDFConfig struct {
	// Page margins
	MarginTop    float64 `json:"margin_top"`
	MarginBottom float64 `json:"margin_bottom"`
	MarginLeft   float64 `json:"margin_left"`
	MarginRight  float64 `json:"margin_right"`

	// Font settings
	FontFamily string  `json:"font_family"`
	FontSize   float64 `json:"font_size"`

	// Colors (RGB values 0-255)
	HeaderColor  [3]int `json:"header_color"`
	ContentColor [3]int `json:"content_color"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Header: HeaderConfig{
			Template: `Some City, {{.date_from}} - {{.date_to}}
Protokół odbioru prac programistycznych

Wykonawca: {{.executor_name}} ({{.executor_email}})
Odbiorca: {{.recipient_name}}

Repozytorium: {{.repository_name}}
- Branch {{.branch_name}}

Lista commitów:`,
			ExecutorName:  "Some Programmer",
			ExecutorEmail: "jan.kowalski@comany.com",
			RecipientName: "CIA",
			Location:      "Warsaw",
		},
		PDF: PDFConfig{
			MarginTop:    20,
			MarginBottom: 20,
			MarginLeft:   20,
			MarginRight:  20,
			FontFamily:   "Arial",
			FontSize:     10,
			HeaderColor:  [3]int{0, 0, 0},
			ContentColor: [3]int{50, 50, 50},
		},
	}
}

// Load loads configuration from file or returns default if file doesn't exist
func Load(configPath string) (*Config, error) {
	// If no config path provided, return default config
	if configPath == "" {
		return DefaultConfig(), nil
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file not found: %s", configPath)
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Parse JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration file: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// Save saves the configuration to a file
func (c *Config) Save(configPath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	return nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Header.Template == "" {
		return fmt.Errorf("header template cannot be empty")
	}

	if c.PDF.FontSize <= 0 {
		return fmt.Errorf("font size must be positive")
	}

	if c.PDF.MarginTop < 0 || c.PDF.MarginBottom < 0 ||
		c.PDF.MarginLeft < 0 || c.PDF.MarginRight < 0 {
		return fmt.Errorf("margins cannot be negative")
	}

	return nil
}

package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App    AppConfig    `yaml:"app"`
	Auth   AuthConfig   `yaml:"auth"`
	HTTP   HTTPConfig   `yaml:"http"`
	AI     AIConfig     `yaml:"ai"`
	Logger LoggerConfig `yaml:"logger"`
}

type AppConfig struct {
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Port        int    `yaml:"port"`
	GinMode     string `yaml:"gin_mode"`
}

type AuthConfig struct {
	APIKey string `yaml:"api_key"`
}

type HTTPConfig struct {
	Timeout          time.Duration `yaml:"timeout"`
	MaxRetries       int           `yaml:"max_retries"`
	RetryWaitTime    time.Duration `yaml:"retry_wait_time"`
	MaxRetryWaitTime time.Duration `yaml:"max_retry_wait_time"`
}

type AIConfig struct {
	BaseURL     string        `yaml:"base_url"`
	APIKey      string        `yaml:"api_key"`
	Model       string        `yaml:"model"`
	Temperature float64       `yaml:"temperature"`
	TopP        float64       `yaml:"top_p"`
	MaxTokens   int           `yaml:"max_tokens"`
	Timeout     time.Duration `yaml:"timeout"`
}

type LoggerConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	expandedData := os.ExpandEnv(string(data))

	var config Config
	if err := yaml.Unmarshal([]byte(expandedData), &config); err != nil {
		return nil, err
	}

	return &config, nil
}

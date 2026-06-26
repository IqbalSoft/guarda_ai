package main

import (
	"encoding/json"
	"os"
)

const configFileName = "config.json"

type AIConfig struct {
	Enabled bool   `json:"enabled"`
	APIKey  string `json:"api_key"`
}

type Config struct {
	Zai      AIConfig `json:"zai"`
	Claude   AIConfig `json:"claude"`
	ChatGPT  AIConfig `json:"chatgpt"`
	Deepseek AIConfig `json:"deepseek"`
	Gemini   AIConfig `json:"gemini"`
	Kimi     AIConfig `json:"kimi"`
}

func loadConfig() Config {
	var cfg Config
	file, err := os.ReadFile(configFileName)
	if err != nil {
		return Config{
			Zai: AIConfig{Enabled: true}, // Default to true so at least one shows
		}
	}
	json.Unmarshal(file, &cfg)
	return cfg
}

func saveConfig(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFileName, data, 0644)
}

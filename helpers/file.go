package helpers

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func ReadTextFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func WriteTextFile(filePath string, content string) error {
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}

type EntryConfig struct {
	Source  string `json:"source"`
	Output  string `json:"output"`
	Lang    string `json:"lang"`
	PkgName string `json:"packageName"`
}

type Config struct {
	Entries []EntryConfig `json:"entries"`
}

func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(filePath string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

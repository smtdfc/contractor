package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Target struct {
	Language string `json:"language"`
	OutDir   string `json:"outDir"`
}

type ContractorConfig struct {
	SourceDir string   `json:"sourceDir"`
	Extension string   `json:"extension"`
	Targets   []Target `json:"targets"`
}

func Load(path string) (*ContractorConfig, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &ContractorConfig{}
	if err := json.Unmarshal(content, cfg); err != nil {
		return nil, fmt.Errorf("invalid config JSON: %w", err)
	}

	if strings.TrimSpace(cfg.SourceDir) == "" {
		return nil, fmt.Errorf("config.sourceDir is required")
	}

	if strings.TrimSpace(cfg.Extension) == "" {
		cfg.Extension = ".contract"
	}

	if !strings.HasPrefix(cfg.Extension, ".") {
		cfg.Extension = "." + cfg.Extension
	}

	if len(cfg.Targets) == 0 {
		return nil, fmt.Errorf("config.targets must contain at least one target")
	}

	for i, target := range cfg.Targets {
		if strings.TrimSpace(target.Language) == "" {
			return nil, fmt.Errorf("config.targets[%d].language is required", i)
		}
		if strings.TrimSpace(target.OutDir) == "" {
			return nil, fmt.Errorf("config.targets[%d].outDir is required", i)
		}
	}

	return cfg, nil
}

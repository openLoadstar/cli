package git

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = "git_config.json"

// Config holds the git remote configuration stored in .loadstar/COMMON/git_config.json.
// NOTE: PAT is stored in plaintext for now. A future version will encrypt this field.
type Config struct {
	RemoteURL string `json:"remote_url"` // e.g. https://github.com/aeolusk/repo.git
	Branch    string `json:"branch"`     // e.g. main
	UserName  string `json:"user_name"`  // git commit author name
	UserEmail string `json:"user_email"` // git commit author email
	PAT       string `json:"pat"`        // personal access token (plaintext, encryption planned)
}

// configPath returns the absolute path to git_config.json.
func configPath(loadstarBase string) string {
	return filepath.Join(loadstarBase, "COMMON", configFileName)
}

// LoadConfig reads git_config.json from .loadstar/COMMON/.
// Returns a zero-value Config (no error) if the file does not exist yet.
func LoadConfig(loadstarBase string) (Config, error) {
	path := configPath(loadstarBase)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, fmt.Errorf("read git config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse git config: %w", err)
	}
	return cfg, nil
}

// SaveConfig writes cfg to .loadstar/COMMON/git_config.json.
func SaveConfig(loadstarBase string, cfg Config) error {
	path := configPath(loadstarBase)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal git config: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// Token returns the PAT. Config file takes priority; falls back to LOADSTAR_GIT_TOKEN env var.
func Token(cfg Config) string {
	if cfg.PAT != "" {
		return cfg.PAT
	}
	return os.Getenv("LOADSTAR_GIT_TOKEN")
}

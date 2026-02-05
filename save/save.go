package save

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// SaveData holds persistent game state.
type SaveData struct {
	HighScore int `json:"high_score"`
	MaxLevel  int `json:"max_level"`
}

func configDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".config", "tankstrike")
}

func savePath() string {
	return filepath.Join(configDir(), "save.json")
}

// Load reads save data from disk.
func Load() *SaveData {
	data, err := os.ReadFile(savePath())
	if err != nil {
		return &SaveData{}
	}
	var s SaveData
	if err := json.Unmarshal(data, &s); err != nil {
		return &SaveData{}
	}
	return &s
}

// Save writes save data to disk.
func Save(s *SaveData) error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(savePath(), data, 0644)
}

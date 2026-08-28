package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/dolfbarr/keychron-battery/internal/model"
)

var defaultCachePath = filepath.Join(os.Getenv("HOME"), ".cache", "keychron", "battery_cache.json")

// Load reads the cached battery states from disk.
func Load() map[string]model.CachedState {
	data, err := os.ReadFile(defaultCachePath)
	if err != nil {
		return make(map[string]model.CachedState)
	}

	var cache map[string]model.CachedState
	if err := json.Unmarshal(data, &cache); err != nil {
		return make(map[string]model.CachedState)
	}
	return cache
}

// Save writes the battery states map to disk.
func Save(cache map[string]model.CachedState) error {
	dir := filepath.Dir(defaultCachePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(defaultCachePath, data, 0644)
}

// Update saves a device's battery state to the cache.
func Update(name string, battery int, charging bool, mode string) {
	cache := Load()
	cache[name] = model.CachedState{
		Battery:   battery,
		Charging:  charging,
		UpdatedAt: time.Now().Unix(),
		Mode:      mode,
	}
	_ = Save(cache)
}

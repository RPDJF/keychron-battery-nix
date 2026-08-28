package cache

import (
	"encoding/json"
	"fmt"
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

// FormatDuration formats seconds into progressive human-readable format.
func FormatDuration(seconds int64) string {
	if seconds < 60 {
		return "Just now"
	}
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24

	if days > 0 {
		remHours := hours % 24
		return fmt.Sprintf("%dd %dh", days, remHours)
	}
	if hours > 0 {
		remMins := minutes % 60
		return fmt.Sprintf("%dh %dm", hours, remMins)
	}
	return fmt.Sprintf("%dm", minutes)
}

// Update saves a device's battery state to the cache.
func Update(name string, battery int, charging bool, mode string) {
	cache := Load()
	now := time.Now().Unix()

	existing, exists := cache[name]
	lastCharged := now
	if exists && existing.LastChargedAt > 0 {
		lastCharged = existing.LastChargedAt
	}

	// If transitioning from charging to discharging, set lastCharged to now
	if exists && existing.Charging && !charging {
		lastCharged = now
	} else if charging {
		lastCharged = now
	}

	cache[name] = model.CachedState{
		Battery:       battery,
		Charging:      charging,
		UpdatedAt:     now,
		LastChargedAt: lastCharged,
		Mode:          mode,
	}
	_ = Save(cache)
}

// GetSinceChargeString computes the formatted duration since last charge.
func GetSinceChargeString(name string, isCharging bool) string {
	if isCharging {
		return "⚡ Charging"
	}

	cache := Load()
	state, exists := cache[name]
	if !exists || state.LastChargedAt == 0 {
		return "󱎫  Active"
	}

	elapsed := time.Now().Unix() - state.LastChargedAt
	if elapsed < 0 {
		elapsed = 0
	}
	return fmt.Sprintf("󱎫  %s", FormatDuration(elapsed))
}

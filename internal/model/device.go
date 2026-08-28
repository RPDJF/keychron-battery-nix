package model

// Device represents a detected Keychron peripheral or receiver.
type Device struct {
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	Type      string `json:"type"`      // "USB", "2.4G", "BT"
	Battery   *int   `json:"battery"`   // Percentage (0-100), nil if unknown
	Charging  bool   `json:"charging"`  // true if charging
	Estimated bool   `json:"estimated"` // true if retrieved from cache
	Signal    string `json:"signal"`    // e.g. "󰤨  100%", "󰒋  Wired", "󰂯  Connected"
}

// CachedState represents the persisted battery state for a device.
type CachedState struct {
	Battery   int   `json:"battery"`
	Charging  bool  `json:"charging"`
	UpdatedAt int64 `json:"updated_at"` // Unix timestamp
	Mode      string `json:"mode"`
}

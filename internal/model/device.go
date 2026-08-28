package model

// DeviceKind represents the peripheral form factor.
type DeviceKind string

const (
	KindKeyboard DeviceKind = "keyboard"
	KindMouse    DeviceKind = "mouse"
	KindGeneric  DeviceKind = "generic"
)

// Device represents a detected Keychron peripheral or receiver.
type Device struct {
	Name        string     `json:"name"`
	Kind        DeviceKind `json:"kind"`
	Icon        string     `json:"icon"`
	Type        string     `json:"type"`         // "USB", "2.4G", "BT"
	Battery     *int       `json:"battery"`      // Percentage (0-100), nil if unknown
	Charging    bool       `json:"charging"`     // true if charging
	Estimated   bool       `json:"estimated"`    // true if retrieved from cache
	SinceCharge string     `json:"since_charge"` // e.g. "3h 15m", "⚡ Charging", "Just now"
	Signal      string     `json:"signal"`       // e.g. "󰤨  100%", "󰒋  Wired", "󰂯  Connected"
	ModelFamily string     `json:"model_family"` // e.g. "max_series", "m_series", "pro_series", "standard_k"
}

// CachedState represents the persisted battery state for a device.
type CachedState struct {
	Battery       int    `json:"battery"`
	Charging      bool   `json:"charging"`
	UpdatedAt     int64  `json:"updated_at"`      // Unix timestamp
	LastChargedAt int64  `json:"last_charged_at"` // Unix timestamp when charging completed / unplugged
	Mode          string `json:"mode"`
	ModelFamily   string `json:"model_family,omitempty"`
}

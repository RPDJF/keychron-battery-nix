package model

// DeviceDriver defines the capabilities of a Keychron model family driver.
type DeviceDriver interface {
	// ID returns the unique identifier for the driver (e.g. "max_series", "pro_series", "m_series").
	ID() string

	// Name returns a human-readable title for the model family.
	Name() string

	// Matches returns true if this driver can handle the given product string, PID, or Bluetooth name.
	Matches(product, pid string) bool

	// Kind returns whether this driver represents a keyboard, mouse, or generic device.
	Kind() DeviceKind

	// Icon returns the primary Nerd Font icon for this device.
	Icon() string

	// ProbeBattery attempts direct low-level feature/raw HID battery interrogation.
	ProbeBattery(nodes []string) (percent int, charging bool, ok bool)

	// SupportsDongle returns true if this model family supports 2.4GHz wireless dongles.
	SupportsDongle() bool
}

// DongleAdapter defines the interface for interrogating specific 2.4GHz USB receivers.
type DongleAdapter interface {
	// ID returns the unique identifier for the dongle protocol (e.g. "keychron_link", "custom_vial").
	ID() string

	// Matches returns true if this adapter handles the given USB receiver product string or PID.
	Matches(product, pid string) bool

	// ProbeCarrier interrogates the dongle for active RF channel & link carrier status.
	ProbeCarrier(nodes []string) (signal string, active bool)
}

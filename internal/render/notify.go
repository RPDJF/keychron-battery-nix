package render

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/dolfbarr/keychron-battery/internal/model"
)

// SendDeviceLowBatteryNotification sends a native FreeDesktop/KDE battery warning for a specific peripheral.
func SendDeviceLowBatteryNotification(dev model.Device) {
	if dev.Battery == nil {
		return
	}

	pct := *dev.Battery
	isMouse := strings.Contains(dev.Name, "M") || strings.Contains(dev.Name, "Mouse")
	icon := "input-keyboard"
	if isMouse {
		icon = "input-mouse"
	}

	title := fmt.Sprintf("%s Battery Low", dev.Name)
	body := fmt.Sprintf("The battery in %s is low (%d%%). Connect it to a charger.", dev.Name, pct)

	urgency := "normal"
	if pct <= 15 {
		urgency = "critical"
	}

	valueHint := fmt.Sprintf("int:value:%d", pct)
	_ = exec.Command("notify-send",
		"-u", urgency,
		"-a", "Power Management",
		"-i", icon,
		"-c", "device",
		"-h", "string:desktop-entry:org.kde.plasma.battery",
		"-h", valueHint,
		title,
		body,
	).Run()
}

// SendStatusSummaryNotification sends an overview notification of all connected peripherals.
func SendStatusSummaryNotification(devices []model.Device) {
	if len(devices) == 0 {
		_ = exec.Command("notify-send", "-a", "Power Management", "Keychron Status", "No Keychron devices detected.").Run()
		return
	}

	var lines []string
	for _, d := range devices {
		bStr := "Active"
		if d.Battery != nil {
			bStr = fmt.Sprintf("%d%%", *d.Battery)
		}
		lines = append(lines, fmt.Sprintf("• %s: %s | %s (%s)", d.Name, bStr, d.Signal, d.Type))
	}

	body := strings.Join(lines, "\n")
	_ = exec.Command("notify-send",
		"-a", "Power Management",
		"-i", "input-keyboard",
		"-c", "device",
		"-h", "string:desktop-entry:org.kde.plasma.battery",
		"Keychron Peripherals",
		body,
	).Run()
}

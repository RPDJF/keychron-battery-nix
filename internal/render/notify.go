package render

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/dolfbarr/keychron-battery/internal/model"
)

// SendNotification triggers a native Linux desktop notification via notify-send.
func SendNotification(devices []model.Device) {
	if len(devices) == 0 {
		_ = exec.Command("notify-send", "-a", "Keychron Battery", "Keychron Status", "No Keychron devices detected.").Run()
		return
	}

	var lines []string
	lowBattery := false
	for _, d := range devices {
		bStr := "Active"
		if d.Battery != nil {
			bStr = fmt.Sprintf("%d%%", *d.Battery)
			if *d.Battery <= 20 {
				lowBattery = true
			}
		}
		lines = append(lines, fmt.Sprintf("• %s: %s | %s (%s)", d.Name, bStr, d.Signal, d.Type))
	}

	body := strings.Join(lines, "\n")
	urgency := "normal"
	title := "Keychron Peripherals"
	if lowBattery {
		urgency = "critical"
		title = "⚠️ Keychron Low Battery!"
	}

	_ = exec.Command("notify-send", "-u", urgency, "-a", "Keychron Battery", "-i", "input-keyboard", title, body).Run()
}

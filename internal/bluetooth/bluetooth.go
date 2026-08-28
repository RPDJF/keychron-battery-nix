package bluetooth

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/dolfbarr/keychron-battery/internal/cache"
	"github.com/dolfbarr/keychron-battery/internal/model"
)

var (
	reModel   = regexp.MustCompile(`(?i)model:\s+(.+)`)
	rePercent = regexp.MustCompile(`(?i)percentage:\s+(\d+)%`)
	reState   = regexp.MustCompile(`(?i)state:\s+(.+)`)
)

// ScanDevices queries UPower and BlueZ for connected Keychron Bluetooth peripherals.
func ScanDevices() []model.Device {
	var devices []model.Device

	out, err := exec.Command("upower", "-e").Output()
	if err != nil {
		return devices
	}

	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.Contains(line, "keyboard") || strings.Contains(line, "mouse") || strings.Contains(strings.ToLower(line), "keychron") {
			infoBytes, err := exec.Command("upower", "-i", line).Output()
			if err != nil {
				continue
			}

			info := string(infoBytes)
			if strings.Contains(strings.ToLower(info), "keychron") || strings.Contains(strings.ToLower(info), "model:") {
				mMatch := reModel.FindStringSubmatch(info)
				pMatch := rePercent.FindStringSubmatch(info)
				sMatch := reState.FindStringSubmatch(info)

				if len(pMatch) > 1 {
					pct, _ := strconv.Atoi(pMatch[1])
					rawName := "Keychron Device"
					if len(mMatch) > 1 {
						rawName = strings.TrimSpace(mMatch[1])
					}

					isMouse := strings.Contains(rawName, "M") || strings.Contains(rawName, "Mouse")
					devName := rawName
					icon := "󰌌 "
					if isMouse {
						icon = "󰍽 "
					}

					isCharging := len(sMatch) > 1 && strings.Contains(strings.ToLower(sMatch[1]), "charging")
					cache.Update(devName, pct, isCharging, "BT")

					devices = append(devices, model.Device{
						Name:        devName,
						Icon:        icon,
						Type:        "󰂯  BT",
						Battery:     &pct,
						Charging:    isCharging,
						Estimated:   false,
						SinceCharge: cache.GetSinceChargeString(devName, isCharging),
						Signal:      "󰂯  Connected",
					})
				}
			}
		}
	}

	return devices
}

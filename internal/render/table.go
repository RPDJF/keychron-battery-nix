package render

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dolfbarr/keychron-battery/internal/model"
)

var ansiRegex = regexp.MustCompile(`\x1B\[[0-?]*[ -/]*[@-~]`)

// VisualLen computes display character cell width taking 2-cell Nerd Font icons into account.
func VisualLen(s string) int {
	clean := ansiRegex.ReplaceAllString(s, "")
	width := 0
	for _, r := range clean {
		if r > 0xE000 {
			width += 2
		} else {
			width++
		}
	}
	return width
}

// PadRight pads string with spaces to reach target visual cell width.
func PadRight(s string, targetWidth int) string {
	v := VisualLen(s)
	if v < targetWidth {
		return s + strings.Repeat(" ", targetWidth-v)
	}
	return s
}

// GetBatteryIcon returns dynamic Nerd Font battery icon.
func GetBatteryIcon(percent int, charging bool) string {
	if charging {
		return "󰂄 "
	}
	switch {
	case percent >= 95:
		return "󰁹 "
	case percent >= 85:
		return "󰂂 "
	case percent >= 75:
		return "󰂁 "
	case percent >= 65:
		return "󰂀 "
	case percent >= 55:
		return "󰁿 "
	case percent >= 45:
		return "󰁾 "
	case percent >= 35:
		return "󰁽 "
	case percent >= 25:
		return "󰁼 "
	case percent >= 15:
		return "󰁻 "
	default:
		return "󰁺 "
	}
}

// GetMiniBar returns a colored mini progress bar.
func GetMiniBar(percent, width int) string {
	filled := (percent * width) / 100
	empty := width - filled
	if empty < 0 {
		empty = 0
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", empty)
}

// RenderBorderless outputs aligned borderless status rows.
func RenderBorderless(devices []model.Device) {
	if len(devices) == 0 {
		fmt.Println("\033[90mNo Keychron devices connected.\033[0m")
		return
	}

	for _, dev := range devices {
		// 1. Device column
		devStyled := fmt.Sprintf("\033[1;37m%s %s\033[0m", dev.Icon, dev.Name)
		devCol := PadRight(devStyled, 24)

		// 2. Battery column
		var battStyled string
		if dev.Battery != nil {
			pct := *dev.Battery
			bar := GetMiniBar(pct, 10)
			bIcon := GetBatteryIcon(pct, dev.Charging)
			bText := fmt.Sprintf("%3d%% %s %s", pct, bIcon, bar)

			color := "\033[1;32m" // Green
			if pct <= 20 {
				color = "\033[1;31m" // Red
			} else if pct <= 50 {
				color = "\033[1;33m" // Yellow
			}
			battStyled = fmt.Sprintf("%s%s\033[0m", color, bText)
		} else {
			battStyled = "\033[32m● Active\033[0m"
		}
		battCol := PadRight(battStyled, 24)

		// 3. Signal column
		sigColor := "\033[32m"
		if strings.Contains(dev.Signal, "Wired") {
			sigColor = "\033[90m"
		}
		sigStyled := fmt.Sprintf("%s%s\033[0m", sigColor, dev.Signal)
		sigCol := PadRight(sigStyled, 14)

		// 4. Mode column
		modeCol := fmt.Sprintf("\033[90m%s\033[0m", dev.Type)

		fmt.Printf("%s%s%s%s\n", devCol, battCol, sigCol, modeCol)
	}
}

// RenderTable outputs formatted table with rounded borders.
func RenderTable(devices []model.Device) {
	fmt.Println("╭────────────────────┬──────────┬─────────┬─────────╮")
	fmt.Println("│ DEVICE             │ BATTERY  │ SIGNAL  │ MODE    │")
	fmt.Println("├────────────────────┼──────────┼─────────┼─────────┤")
	for _, dev := range devices {
		devStr := fmt.Sprintf("%s %s", dev.Icon, dev.Name)
		battStr := "● Active"
		if dev.Battery != nil {
			pct := *dev.Battery
			bar := GetMiniBar(pct, 8)
			bIcon := GetBatteryIcon(pct, dev.Charging)
			battStr = fmt.Sprintf("%3d%% %s%s", pct, bIcon, bar)
		}
		fmt.Printf("│ %-18s │ %-8s │ %-7s │ %-7s │\n", devStr, battStr, dev.Signal, dev.Type)
	}
	fmt.Println("╰────────────────────┴──────────┴─────────┴─────────╯")
}

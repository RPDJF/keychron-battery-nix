package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/dolfbarr/keychron-battery/internal/bluetooth"
	"github.com/dolfbarr/keychron-battery/internal/daemon"
	"github.com/dolfbarr/keychron-battery/internal/hid"
	"github.com/dolfbarr/keychron-battery/internal/model"
	"github.com/dolfbarr/keychron-battery/internal/render"
)

var (
	Version   = "1.0.0"
	BuildDate = "2026-08-28"
)

func getCombinedDevices() []model.Device {
	devices := hid.DetectDevices()
	btDevices := bluetooth.ScanDevices()

	for _, bt := range btDevices {
		found := false
		for _, d := range devices {
			if d.Name == bt.Name {
				found = true
				break
			}
		}
		if !found {
			devices = append(devices, bt)
		}
	}
	return devices
}

func main() {
	jsonFlag := flag.Bool("json", false, "Output results in JSON format")
	notifyFlag := flag.Bool("notify", false, "Send desktop notification via notify-send")
	watchFlag := flag.Int("watch", 0, "Continuous watch mode with refresh interval in seconds")
	tableFlag := flag.Bool("table", false, "Render as a table with rounded borders")
	daemonFlag := flag.Bool("daemon", false, "Run as background monitoring daemon")
	intervalFlag := flag.Int("interval", 30, "Interval in seconds for daemon mode")
	versionFlag := flag.Bool("version", false, "Show version and exit")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("keychron-battery v%s (%s)\n", Version, BuildDate)
		return
	}

	if *daemonFlag {
		daemon.Run(time.Duration(*intervalFlag) * time.Second)
		return
	}

	for {
		devices := getCombinedDevices()

		if *notifyFlag {
			render.SendNotification(devices)
			if *watchFlag == 0 {
				break
			}
		}

		if *jsonFlag {
			data, _ := json.MarshalIndent(devices, "", "  ")
			fmt.Println(string(data))
		} else if *tableFlag {
			render.RenderTable(devices)
		} else {
			render.RenderBorderless(devices)
		}

		if *watchFlag > 0 {
			time.Sleep(time.Duration(*watchFlag) * time.Second)
			if !*jsonFlag {
				fmt.Print("\033[H\033[J")
			}
		} else {
			break
		}
	}
}

package daemon

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dolfbarr/keychron-battery/internal/bluetooth"
	"github.com/dolfbarr/keychron-battery/internal/hid"
	"github.com/dolfbarr/keychron-battery/internal/model"
	"github.com/dolfbarr/keychron-battery/internal/render"
)

// Run starts the background monitoring daemon.
func Run(interval time.Duration) {
	log.Printf("Starting Keychron Battery Daemon (interval: %v)...", interval)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Initial scan
	pollAndNotify()

	for {
		select {
		case <-sigChan:
			log.Println("Received shutdown signal. Exiting daemon...")
			return
		case <-ticker.C:
			pollAndNotify()
		case <-ctx.Done():
			return
		}
	}
}

var lastLowBatteryNotified = false

func pollAndNotify() {
	devices := hid.DetectDevices()
	btDevices := bluetooth.ScanDevices()

	// Deduplicate Bluetooth vs USB
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

	hasLowBattery := false
	for _, d := range devices {
		if d.Battery != nil && *d.Battery <= 20 && !d.Charging {
			hasLowBattery = true
			break
		}
	}

	if hasLowBattery && !lastLowBatteryNotified {
		render.SendNotification(devices)
		lastLowBatteryNotified = true
	} else if !hasLowBattery {
		lastLowBatteryNotified = false
	}
}

// TriggerTestAlert executes the exact daemon notification pipeline with a simulated low-battery event.
func TriggerTestAlert() {
	lowPct := 15
	fullPct := 100
	mockDevices := []model.Device{
		{
			Name:      "Keychron M6",
			Icon:      "󰍽 ",
			Type:      "󰖩  2.4G",
			Battery:   &lowPct,
			Charging:  false,
			Estimated: true,
			Signal:    "󰤨  100%",
		},
		{
			Name:      "Keychron K3 Max",
			Icon:      "󰌌 ",
			Type:      "󰖩  2.4G",
			Battery:   &fullPct,
			Charging:  false,
			Estimated: true,
			Signal:    "󰤨  100%",
		},
	}
	render.SendNotification(mockDevices)
}

package hid

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dolfbarr/keychron-battery/internal/cache"
	"github.com/dolfbarr/keychron-battery/internal/dongles"
	"github.com/dolfbarr/keychron-battery/internal/drivers"
	"github.com/dolfbarr/keychron-battery/internal/model"
)

// USBDevice represents a top-level Keychron USB device.
type USBDevice struct {
	SysfsPath string
	Product   string
	ProductID string
	Nodes     []string
}

// GetParentUSBPath resolves the parent USB sysfs directory for a hidraw path.
func GetParentUSBPath(hidrawPath string) string {
	devLink := filepath.Join(hidrawPath, "device")
	target, err := filepath.EvalSymlinks(devLink)
	if err != nil {
		return ""
	}

	for target != "/" && strings.Contains(filepath.Base(target), ":") {
		target = filepath.Dir(target)
	}
	return target
}

// ScanUSBDevices discovers all connected Keychron USB devices and maps their child hidraw nodes.
func ScanUSBDevices() []USBDevice {
	var devices []USBDevice

	matches, _ := filepath.Glob("/sys/bus/usb/devices/*/uevent")
	for _, uevent := range matches {
		dir := filepath.Dir(uevent)
		base := filepath.Base(dir)
		if strings.Contains(base, ":") || strings.HasPrefix(base, "usb") {
			continue
		}

		data, err := os.ReadFile(uevent)
		if err != nil {
			continue
		}

		content := string(data)
		if strings.Contains(content, "PRODUCT=3434/") || strings.Contains(content, "VENDOR=3434") {
			product := "Keychron Device"
			if pData, err := os.ReadFile(filepath.Join(dir, "product")); err == nil {
				product = strings.TrimSpace(string(pData))
			}

			prodID := ""
			if idData, err := os.ReadFile(filepath.Join(dir, "idProduct")); err == nil {
				prodID = strings.ToLower(strings.TrimSpace(string(idData)))
			}

			devices = append(devices, USBDevice{
				SysfsPath: dir,
				Product:   product,
				ProductID: prodID,
			})
		}
	}

	hidrawMatches, _ := filepath.Glob("/sys/class/hidraw/hidraw*")
	for _, h := range hidrawMatches {
		parent := GetParentUSBPath(h)
		if parent == "" {
			continue
		}

		for i := range devices {
			if devices[i].SysfsPath == parent {
				nodeName := filepath.Base(h)
				devices[i].Nodes = append(devices[i].Nodes, "/dev/"+nodeName)
				break
			}
		}
	}

	return devices
}

// DetectDevices scans USB devices and delegates decoding to pluggable drivers & dongle adapters.
func DetectDevices() []model.Device {
	var results []model.Device
	usbDevs := ScanUSBDevices()
	cachedStates := cache.Load()

	hasWiredMouse := false
	hasWiredKb := false

	cachedKbName := "Keychron K3 Max"
	cachedMouseName := "Keychron M6"
	for name, state := range cachedStates {
		if strings.Contains(name, "M") || strings.Contains(name, "Mouse") || state.ModelFamily == "m_series" {
			cachedMouseName = name
		} else if strings.Contains(name, "K") || strings.Contains(name, "Q") || strings.Contains(name, "V") || strings.Contains(name, "Lemokey") {
			cachedKbName = name
		}
	}

	// 1. Process direct USB connected devices using drivers
	for _, dev := range usbDevs {
		dongleAdapter := dongles.FindAdapter(dev.Product, dev.ProductID)
		if dongleAdapter != nil {
			continue
		}

		driver := drivers.FindDriver(dev.Product, dev.ProductID)
		devName := dev.Product
		if devName == "" || devName == "Keychron Device" {
			if driver.Kind() == model.KindMouse {
				devName = cachedMouseName
			} else {
				devName = cachedKbName
			}
		}

		if driver.Kind() == model.KindMouse {
			hasWiredMouse = true
			cachedMouseName = devName
		} else if driver.Kind() == model.KindKeyboard {
			hasWiredKb = true
			cachedKbName = devName
		}

		battPct := 100
		isCharging := true
		if pct, chg, ok := driver.ProbeBattery(dev.Nodes); ok {
			battPct = pct
			isCharging = chg
		}

		cache.Update(devName, battPct, isCharging, "USB")

		results = append(results, model.Device{
			Name:        devName,
			Kind:        driver.Kind(),
			Icon:        driver.Icon(),
			Type:        "󰒋  USB",
			Battery:     &battPct,
			Charging:    isCharging,
			Estimated:   false,
			SinceCharge: cache.GetSinceChargeString(devName, isCharging),
			Signal:      "󰒋  Wired",
			ModelFamily: driver.ID(),
		})
	}

	// 2. Process 2.4GHz Dongles using pluggable dongle adapters
	dongleCount := 0
	for _, dev := range usbDevs {
		dongleAdapter := dongles.FindAdapter(dev.Product, dev.ProductID)
		if dongleAdapter == nil {
			continue
		}

		dongleCount++
		signalStr, _ := dongleAdapter.ProbeCarrier(dev.Nodes)

		if dongleCount == 1 && !hasWiredKb {
			kbDriver := drivers.FindDriver(cachedKbName, "")
			devName := cachedKbName
			var battPtr *int
			isEst := false
			if c, ok := cachedStates[devName]; ok {
				val := c.Battery
				battPtr = &val
				isEst = true
			} else {
				val := 100
				battPtr = &val
			}

			results = append(results, model.Device{
				Name:        devName,
				Kind:        model.KindKeyboard,
				Icon:        kbDriver.Icon(),
				Type:        "󰖩  2.4G",
				Battery:     battPtr,
				Charging:    false,
				Estimated:   isEst,
				SinceCharge: cache.GetSinceChargeString(devName, false),
				Signal:      signalStr,
				ModelFamily: kbDriver.ID(),
			})
		} else if dongleCount == 2 && !hasWiredMouse {
			mouseDriver := drivers.FindDriver(cachedMouseName, "")
			devName := cachedMouseName
			var battPtr *int
			isEst := false
			if c, ok := cachedStates[devName]; ok {
				val := c.Battery
				battPtr = &val
				isEst = true
			} else {
				val := 98
				battPtr = &val
			}

			results = append(results, model.Device{
				Name:        devName,
				Kind:        model.KindMouse,
				Icon:        mouseDriver.Icon(),
				Type:        "󰖩  2.4G",
				Battery:     battPtr,
				Charging:    false,
				Estimated:   isEst,
				SinceCharge: cache.GetSinceChargeString(devName, false),
				Signal:      signalStr,
				ModelFamily: mouseDriver.ID(),
			})
		}
	}

	return results
}

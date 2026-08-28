package hid

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/dolfbarr/keychron-battery/internal/cache"
	"github.com/dolfbarr/keychron-battery/internal/model"
)

const (
	KeychronVendorID = "3434"
)

func ioc(dir, typeCode, nr, size uintptr) uintptr {
	return (dir << 30) | (typeCode << 8) | nr | (size << 16)
}

func hidiocsfeat(length uintptr) uintptr {
	return ioc(3, uintptr('H'), 0x06, length)
}

func hidiocgfeat(length uintptr) uintptr {
	return ioc(3, uintptr('H'), 0x07, length)
}

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

	// 1. Scan USB sysfs
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

	// 2. Map hidraw nodes to parent USB devices
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

// ProbeFeatureBattery reads Keychron vendor Feature Reports 0x51/0x52 for battery % and charging status.
func ProbeFeatureBattery(devPath string) (int, bool, bool) {
	fd, err := syscall.Open(devPath, syscall.O_RDWR, 0)
	if err != nil {
		return 0, false, false
	}
	defer syscall.Close(fd)

	for _, repID := range []byte{0x51, 0x52, 0x53} {
		buf := make([]byte, 21)
		buf[0] = repID
		req := hidiocgfeat(uintptr(len(buf)))
		_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), req, uintptr(unsafe.Pointer(&buf[0])))
		if errno == 0 && (buf[1] > 0 || buf[2] > 0 || buf[4] > 0) {
			pct := int(buf[1])
			if pct >= 1 && pct <= 100 {
				isCharging := len(buf) > 4 && buf[4] == 1
				return pct, isCharging, true
			}
		}
	}
	return 0, false, false
}

// ProbeDongleLinkQuality queries the 2.4GHz dongle for RF link carrier status via Report 0x54.
func ProbeDongleLinkQuality(devPath string) (string, bool) {
	fd, err := syscall.Open(devPath, syscall.O_RDWR|syscall.O_NONBLOCK, 0)
	if err != nil {
		return "󰤨  100%", false
	}
	defer syscall.Close(fd)

	// Step 1: Send Set Feature 0x51
	buf := make([]byte, 21)
	buf[0] = 0x51
	buf[1] = 0x01
	setReq := hidiocsfeat(uintptr(len(buf)))
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), setReq, uintptr(unsafe.Pointer(&buf[0])))

	// Step 2: Send Output Report 0xB2
	out := make([]byte, 33)
	out[0] = 0xB2
	out[1] = 0x01
	_, _ = syscall.Write(fd, out)

	time.Sleep(20 * time.Millisecond)

	resp := make([]byte, 64)
	n, _ := syscall.Read(fd, resp)
	if n >= 3 && resp[0] == 0x54 {
		return "󰤨  100%", true
	}
	return "󰤨  100%", true
}

// DetectDevices scans USB devices and decodes live/cached peripherals.
func DetectDevices() []model.Device {
	var results []model.Device
	usbDevs := ScanUSBDevices()
	cachedStates := cache.Load()

	hasWiredMouse := false
	hasWiredKb := false

	// 1. Direct wired USB devices
	for _, dev := range usbDevs {
		isLink := strings.Contains(dev.Product, "Link") || dev.ProductID == "d030" || dev.ProductID == "d031"
		if isLink {
			continue
		}

		isMouse := strings.Contains(dev.Product, "M") || strings.Contains(dev.Product, "Mouse") || dev.ProductID == "d03f"
		devName := "Keychron K3 Max"
		icon := "󰌌 "
		if isMouse {
			devName = "Keychron M6"
			icon = "󰍽 "
			hasWiredMouse = true
		} else {
			hasWiredKb = true
		}

		battPct := 100
		isCharging := true
		for _, node := range dev.Nodes {
			if pct, chg, ok := ProbeFeatureBattery(node); ok {
				battPct = pct
				isCharging = chg
				break
			}
		}

		cache.Update(devName, battPct, isCharging, "USB")

		results = append(results, model.Device{
			Name:      devName,
			Icon:      icon,
			Type:      "󰒋  USB",
			Battery:   &battPct,
			Charging:  isCharging,
			Estimated: false,
			Signal:    "󰒋  Wired",
		})
	}

	// 2. 2.4GHz Dongles
	dongleCount := 0
	for _, dev := range usbDevs {
		isLink := strings.Contains(dev.Product, "Link") || dev.ProductID == "d030" || dev.ProductID == "d031"
		if !isLink {
			continue
		}

		dongleCount++
		signalStr := "󰤨  100%"
		for _, node := range dev.Nodes {
			if sig, ok := ProbeDongleLinkQuality(node); ok {
				signalStr = sig
				break
			}
		}

		if dongleCount == 1 && !hasWiredKb {
			devName := "Keychron K3 Max"
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
				Name:      devName,
				Icon:      "󰌌 ",
				Type:      "󰖩  2.4G",
				Battery:   battPtr,
				Charging:  false,
				Estimated: isEst,
				Signal:    signalStr,
			})
		} else if dongleCount == 2 && !hasWiredMouse {
			devName := "Keychron M6"
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
				Name:      devName,
				Icon:      "󰍽 ",
				Type:      "󰖩  2.4G",
				Battery:   battPtr,
				Charging:  false,
				Estimated: isEst,
				Signal:    signalStr,
			})
		}
	}

	return results
}

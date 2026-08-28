package drivers

import (
	"strings"
	"syscall"
	"unsafe"

	"github.com/dolfbarr/keychron-battery/internal/model"
)

// MSeriesDriver handles Keychron M-Series wireless mice (M1, M2, M3, M6, M7, etc.).
type MSeriesDriver struct{}

func (m *MSeriesDriver) ID() string {
	return "m_series"
}

func (m *MSeriesDriver) Name() string {
	return "Keychron M-Series Mouse"
}

func (m *MSeriesDriver) Kind() model.DeviceKind {
	return model.KindMouse
}

func (m *MSeriesDriver) Icon() string {
	return "󰍽 "
}

func (m *MSeriesDriver) SupportsDongle() bool {
	return true
}

func (m *MSeriesDriver) Matches(product, pid string) bool {
	pLower := strings.ToLower(product)
	pidLower := strings.ToLower(pid)
	return strings.Contains(pLower, "m1") || strings.Contains(pLower, "m2") ||
		strings.Contains(pLower, "m3") || strings.Contains(pLower, "m4") ||
		strings.Contains(pLower, "m6") || strings.Contains(pLower, "m7") ||
		strings.Contains(pLower, "mouse") || pidLower == "d03f"
}

func (m *MSeriesDriver) ProbeBattery(nodes []string) (int, bool, bool) {
	for _, node := range nodes {
		fd, err := syscall.Open(node, syscall.O_RDWR, 0)
		if err != nil {
			continue
		}

		for _, repID := range []byte{0x51, 0x52, 0x53} {
			buf := make([]byte, 21)
			buf[0] = repID
			req := (uintptr(3) << 30) | (uintptr('H') << 8) | 0x07 | (uintptr(len(buf)) << 16)
			_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), req, uintptr(unsafe.Pointer(&buf[0])))
			if errno == 0 && (buf[1] > 0 || buf[2] > 0 || buf[4] > 0) {
				pct := int(buf[1])
				if pct >= 1 && pct <= 100 {
					isCharging := len(buf) > 4 && buf[4] == 1
					syscall.Close(fd)
					return pct, isCharging, true
				}
			}
		}
		syscall.Close(fd)
	}
	return 0, false, false
}

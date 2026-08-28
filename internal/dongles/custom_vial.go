package dongles

import (
	"strings"
	"syscall"
	"time"
)

// CustomVialAdapter provides extensible support for DIY/Vial 2.4GHz USB wireless dongles.
type CustomVialAdapter struct{}

func (c *CustomVialAdapter) ID() string {
	return "custom_vial"
}

func (c *CustomVialAdapter) Matches(product, pid string) bool {
	pLower := strings.ToLower(product)
	return strings.Contains(pLower, "vial") || strings.Contains(pLower, "qmk dongle") || strings.Contains(pLower, "nrf52")
}

func (c *CustomVialAdapter) ProbeCarrier(nodes []string) (string, bool) {
	for _, node := range nodes {
		fd, err := syscall.Open(node, syscall.O_RDWR|syscall.O_NONBLOCK, 0)
		if err != nil {
			continue
		}

		// Probe standard VIA protocol version on Usage Page 0xFF60
		pkt := make([]byte, 32)
		pkt[0] = 0x01 // VIA get protocol version
		_, _ = syscall.Write(fd, pkt)

		time.Sleep(20 * time.Millisecond)

		resp := make([]byte, 32)
		n, _ := syscall.Read(fd, resp)
		syscall.Close(fd)

		if n > 0 && resp[0] != 0xFF {
			return "󰤨  100%", true
		}
	}
	return "󰤨  100%", true
}

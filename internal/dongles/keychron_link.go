package dongles

import (
	"strings"
	"syscall"
	"time"
	"unsafe"
)

// KeychronLinkAdapter handles official Keychron 2.4GHz Link Receivers (3434:d030 / 3434:d031).
type KeychronLinkAdapter struct{}

func (k *KeychronLinkAdapter) ID() string {
	return "keychron_link"
}

func (k *KeychronLinkAdapter) Matches(product, pid string) bool {
	pLower := strings.ToLower(product)
	pidLower := strings.ToLower(pid)
	return strings.Contains(pLower, "link") || pidLower == "d030" || pidLower == "d031"
}

func (k *KeychronLinkAdapter) ProbeCarrier(nodes []string) (string, bool) {
	for _, node := range nodes {
		fd, err := syscall.Open(node, syscall.O_RDWR|syscall.O_NONBLOCK, 0)
		if err != nil {
			continue
		}

		// Step 1: Set Feature Report 0x51
		buf := make([]byte, 21)
		buf[0] = 0x51
		buf[1] = 0x01
		setReq := (uintptr(3) << 30) | (uintptr('H') << 8) | 0x06 | (uintptr(len(buf)) << 16)
		_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), setReq, uintptr(unsafe.Pointer(&buf[0])))

		// Step 2: Send Output Report 0xB2
		out := make([]byte, 33)
		out[0] = 0xB2
		out[1] = 0x01
		_, _ = syscall.Write(fd, out)

		time.Sleep(20 * time.Millisecond)

		resp := make([]byte, 64)
		n, _ := syscall.Read(fd, resp)
		syscall.Close(fd)

		if n >= 3 && resp[0] == 0x54 {
			return "󰤨  100%", true
		}
	}
	return "󰤨  100%", true
}

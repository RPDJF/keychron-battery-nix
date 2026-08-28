package drivers

import (
	"strings"

	"github.com/dolfbarr/keychron-battery/internal/model"
)

// StandardKSeriesDriver handles classic Keychron K-series keyboards (K2, K3, K4, K6, K8, K10, etc.).
// Designed with extensible hooks for future official or DIY 2.4GHz wireless dongles.
type StandardKSeriesDriver struct{}

func (s *StandardKSeriesDriver) ID() string {
	return "standard_k"
}

func (s *StandardKSeriesDriver) Name() string {
	return "Keychron K-Series Keyboard"
}

func (s *StandardKSeriesDriver) Kind() model.DeviceKind {
	return model.KindKeyboard
}

func (s *StandardKSeriesDriver) Icon() string {
	return "󰌌 "
}

func (s *StandardKSeriesDriver) SupportsDongle() bool {
	// Extensible: can be paired with custom or future wireless dongles
	return true
}

func (s *StandardKSeriesDriver) Matches(product, pid string) bool {
	pLower := strings.ToLower(product)
	return strings.HasPrefix(pLower, "k") || strings.Contains(pLower, "k2") ||
		strings.Contains(pLower, "k3") || strings.Contains(pLower, "k4") ||
		strings.Contains(pLower, "k6") || strings.Contains(pLower, "k8")
}

func (s *StandardKSeriesDriver) ProbeBattery(nodes []string) (int, bool, bool) {
	return 100, true, true
}

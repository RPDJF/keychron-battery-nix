package drivers

import (
	"strings"

	"github.com/dolfbarr/keychron-battery/internal/model"
)

// ProSeriesDriver handles Keychron Pro keyboards (K Pro, Q Pro, V Pro).
// Designed with extensible hooks for future official or DIY/Vial 2.4GHz dongles.
type ProSeriesDriver struct{}

func (p *ProSeriesDriver) ID() string {
	return "pro_series"
}

func (p *ProSeriesDriver) Name() string {
	return "Keychron Pro Series Keyboard"
}

func (p *ProSeriesDriver) Kind() model.DeviceKind {
	return model.KindKeyboard
}

func (p *ProSeriesDriver) Icon() string {
	return "󰌌 "
}

func (p *ProSeriesDriver) SupportsDongle() bool {
	// Extensible: can be paired with custom / future dongles
	return true
}

func (p *ProSeriesDriver) Matches(product, pid string) bool {
	pLower := strings.ToLower(product)
	return strings.Contains(pLower, "pro") && (strings.Contains(pLower, "k") || strings.Contains(pLower, "q") || strings.Contains(pLower, "v"))
}

func (p *ProSeriesDriver) ProbeBattery(nodes []string) (int, bool, bool) {
	return 100, true, true
}

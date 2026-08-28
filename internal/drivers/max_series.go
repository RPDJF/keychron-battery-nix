package drivers

import (
	"strings"

	"github.com/dolfbarr/keychron-battery/internal/model"
)

// MaxSeriesDriver handles Keychron Max series wireless keyboards (K Max, Q Max, V Max, Lemokey).
type MaxSeriesDriver struct{}

func (k *MaxSeriesDriver) ID() string {
	return "max_series"
}

func (k *MaxSeriesDriver) Name() string {
	return "Keychron Max Series Keyboard"
}

func (k *MaxSeriesDriver) Kind() model.DeviceKind {
	return model.KindKeyboard
}

func (k *MaxSeriesDriver) Icon() string {
	return "󰌌 "
}

func (k *MaxSeriesDriver) SupportsDongle() bool {
	return true
}

func (k *MaxSeriesDriver) Matches(product, pid string) bool {
	pLower := strings.ToLower(product)
	return strings.Contains(pLower, "max") || strings.Contains(pLower, "lemokey")
}

func (k *MaxSeriesDriver) ProbeBattery(nodes []string) (int, bool, bool) {
	// Keyboards currently charge when on direct USB cable
	return 100, true, true
}

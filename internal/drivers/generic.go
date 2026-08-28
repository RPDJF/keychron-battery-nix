package drivers

import (
	"strings"

	"github.com/dolfbarr/keychron-battery/internal/model"
)

// GenericDriver handles any unrecognized Keychron device.
type GenericDriver struct{}

func (g *GenericDriver) ID() string {
	return "generic"
}

func (g *GenericDriver) Name() string {
	return "Keychron Peripheral"
}

func (g *GenericDriver) Kind() model.DeviceKind {
	return model.KindGeneric
}

func (g *GenericDriver) Icon() string {
	return "🔌 "
}

func (g *GenericDriver) SupportsDongle() bool {
	return true
}

func (g *GenericDriver) Matches(product, pid string) bool {
	return true
}

func (g *GenericDriver) ProbeBattery(nodes []string) (int, bool, bool) {
	isMouse := false
	for _, n := range nodes {
		if strings.Contains(strings.ToLower(n), "mouse") {
			isMouse = true
			break
		}
	}
	if isMouse {
		m := &MSeriesDriver{}
		return m.ProbeBattery(nodes)
	}
	return 100, true, true
}

package drivers

import (
	"github.com/dolfbarr/keychron-battery/internal/model"
)

var registeredDrivers []model.DeviceDriver

func init() {
	Register(&MSeriesDriver{})
	Register(&MaxSeriesDriver{})
	Register(&ProSeriesDriver{})
	Register(&StandardKSeriesDriver{})
	Register(&GenericDriver{})
}

// Register adds a new device driver to the global registry.
func Register(driver model.DeviceDriver) {
	registeredDrivers = append(registeredDrivers, driver)
}

// FindDriver returns the most specific driver matching the given product and PID.
func FindDriver(product, pid string) model.DeviceDriver {
	for _, d := range registeredDrivers {
		if d.Matches(product, pid) {
			return d
		}
	}
	return &GenericDriver{}
}

// FindDriverByKind returns the best driver matching the given device kind.
func FindDriverByKind(kind model.DeviceKind) model.DeviceDriver {
	for _, d := range registeredDrivers {
		if d.Kind() == kind {
			return d
		}
	}
	return &GenericDriver{}
}

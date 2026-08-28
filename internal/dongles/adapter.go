package dongles

import (
	"github.com/dolfbarr/keychron-battery/internal/model"
)

var registeredAdapters []model.DongleAdapter

func init() {
	Register(&KeychronLinkAdapter{})
	Register(&CustomVialAdapter{})
}

// Register adds a new 2.4GHz dongle protocol adapter.
func Register(adapter model.DongleAdapter) {
	registeredAdapters = append(registeredAdapters, adapter)
}

// FindAdapter searches for a matching dongle adapter for the given USB product string and PID.
func FindAdapter(product, pid string) model.DongleAdapter {
	for _, a := range registeredAdapters {
		if a.Matches(product, pid) {
			return a
		}
	}
	return nil
}

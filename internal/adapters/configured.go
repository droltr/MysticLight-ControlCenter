package adapters

import (
	"os"
	"time"

	"github.com/droltr/MysticLight-ControlCenter/internal/adapters/coolercontrol"
	"github.com/droltr/MysticLight-ControlCenter/internal/adapters/openrgb"
	"github.com/droltr/MysticLight-ControlCenter/internal/provider"
)

func FromEnvironment() []provider.Adapter {
	adapters := DefaultReadOnlyAdapters()
	if address := os.Getenv("MYSTICLIGHT_OPENRGB_ADDRESS"); address != "" {
		adapters = replace(adapters, OpenRGB, openrgb.New(address, 3*time.Second))
	}
	if baseURL := os.Getenv("MYSTICLIGHT_COOLERCONTROL_URL"); baseURL != "" {
		adapters = replace(adapters, CoolerControl, coolercontrol.New(
			baseURL,
			os.Getenv("MYSTICLIGHT_COOLERCONTROL_PATH"),
			os.Getenv("MYSTICLIGHT_COOLERCONTROL_TOKEN"),
			3*time.Second,
		))
	}
	return adapters
}

func replace(adapters []provider.Adapter, name string, replacement provider.Adapter) []provider.Adapter {
	for index, adapter := range adapters {
		if adapter.Name() == name {
			adapters[index] = replacement
			return adapters
		}
	}
	return append(adapters, replacement)
}

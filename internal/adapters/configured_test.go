package adapters

import "testing"

func TestFromEnvironmentKeepsProvidersUnconfiguredByDefault(t *testing.T) {
	t.Setenv("MYSTICLIGHT_OPENRGB_ADDRESS", "")
	t.Setenv("MYSTICLIGHT_COOLERCONTROL_URL", "")
	for _, adapter := range FromEnvironment() {
		if adapter.Name() == OpenRGB || adapter.Name() == CoolerControl {
			if adapter.Health(nil).Available {
				t.Fatalf("provider %q should be unavailable without configuration", adapter.Name())
			}
		}
	}
}

func TestFromEnvironmentConfiguresReadOnlyProviders(t *testing.T) {
	t.Setenv("MYSTICLIGHT_OPENRGB_ADDRESS", "127.0.0.1:6742")
	t.Setenv("MYSTICLIGHT_COOLERCONTROL_URL", "http://127.0.0.1:11987")
	t.Setenv("MYSTICLIGHT_COOLERCONTROL_PATH", "/read-only")
	t.Setenv("MYSTICLIGHT_COOLERCONTROL_TOKEN", "test-only-token")
	for _, adapter := range FromEnvironment() {
		if adapter.Name() == OpenRGB || adapter.Name() == CoolerControl {
			if len(adapter.Capabilities()) != 1 || adapter.Capabilities()[0].Supports("apply") {
				t.Fatalf("provider %q is not read-only", adapter.Name())
			}
		}
	}
}

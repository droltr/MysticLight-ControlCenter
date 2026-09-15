package profile

import (
	"testing"

	"github.com/droltr/MysticLight-ControlCenter/internal/domain"
)

func TestRuntimeActivatesValidatedProfile(t *testing.T) {
	providers := map[string]struct{}{"openrgb": {}, "coolercontrol": {}}
	fallback := domain.Profile{
		Name:      "safe",
		Ownership: []domain.Ownership{{Resource: domain.ResourceRGB, Provider: "openrgb"}},
	}
	runtime, err := NewRuntime(providers, fallback)
	if err != nil {
		t.Fatalf("failed to create runtime: %v", err)
	}
	profile := domain.Profile{
		Name: "quiet",
		Ownership: []domain.Ownership{
			{Resource: domain.ResourceRGB, Provider: "openrgb"},
			{Resource: domain.ResourcePWM, Provider: "coolercontrol"},
		},
	}
	if err := runtime.Activate(profile); err != nil {
		t.Fatalf("failed to activate profile: %v", err)
	}
	if runtime.Active().Name != "quiet" {
		t.Fatalf("unexpected active profile: %q", runtime.Active().Name)
	}
}

func TestRuntimeRejectsInvalidProfileAndRecoversFallback(t *testing.T) {
	providers := map[string]struct{}{"openrgb": {}}
	fallback := domain.Profile{
		Name:      "safe",
		Ownership: []domain.Ownership{{Resource: domain.ResourceRGB, Provider: "openrgb"}},
	}
	runtime, err := NewRuntime(providers, fallback)
	if err != nil {
		t.Fatalf("failed to create runtime: %v", err)
	}
	err = runtime.Activate(domain.Profile{
		Name:      "unsafe",
		Ownership: []domain.Ownership{{Resource: domain.ResourcePWM, Provider: "missing"}},
	})
	if err == nil {
		t.Fatal("expected invalid profile activation to fail")
	}
	if recovered := runtime.Recover(); recovered.Name != "safe" {
		t.Fatalf("unexpected recovery profile: %q", recovered.Name)
	}
}

func TestRuntimeActivatesJSONProfile(t *testing.T) {
	runtime, err := NewRuntime(
		map[string]struct{}{"openrgb": {}},
		domain.Profile{Name: "safe", Ownership: []domain.Ownership{{Resource: domain.ResourceRGB, Provider: "openrgb"}}},
	)
	if err != nil {
		t.Fatalf("failed to create runtime: %v", err)
	}
	if err := runtime.ActivateJSON([]byte(`{"name":"quiet","ownership":[{"resource":"rgb","provider":"openrgb"}]}`)); err != nil {
		t.Fatalf("failed to activate JSON profile: %v", err)
	}
	if runtime.Active().Name != "quiet" {
		t.Fatalf("unexpected active profile: %q", runtime.Active().Name)
	}
}

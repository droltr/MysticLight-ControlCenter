package domain

import "testing"

func TestProfileValidateRejectsConflictingOwners(t *testing.T) {
	profile := Profile{
		Name: "quiet",
		Ownership: []Ownership{
			{Resource: ResourceRGB, Provider: "openrgb"},
			{Resource: ResourceRGB, Provider: "other"},
		},
	}

	err := profile.Validate(map[string]struct{}{"openrgb": {}, "other": {}})
	if err == nil {
		t.Fatal("expected conflicting ownership to be rejected")
	}
}

func TestProfileValidateRejectsUnknownProvider(t *testing.T) {
	profile := Profile{
		Name:      "quiet",
		Ownership: []Ownership{{Resource: ResourceRGB, Provider: "missing"}},
	}

	err := profile.Validate(map[string]struct{}{"openrgb": {}})
	if err == nil {
		t.Fatal("expected unknown provider to be rejected")
	}
}

func TestProfileValidateAcceptsUniqueKnownOwners(t *testing.T) {
	profile := Profile{
		Name: "quiet",
		Ownership: []Ownership{
			{Resource: ResourceRGB, Provider: "openrgb"},
			{Resource: ResourcePWM, Provider: "coolercontrol"},
		},
	}

	if err := profile.Validate(map[string]struct{}{"openrgb": {}, "coolercontrol": {}}); err != nil {
		t.Fatalf("expected profile to validate: %v", err)
	}
}

func TestProviderHealthState(t *testing.T) {
	cases := []struct {
		name     string
		health   ProviderHealth
		expected HealthState
	}{
		{name: "available", health: ProviderHealth{Available: true}, expected: HealthAvailable},
		{name: "degraded", health: ProviderHealth{Available: true, Error: "timeout"}, expected: HealthDegraded},
		{name: "unavailable", health: ProviderHealth{}, expected: HealthUnavailable},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if actual := test.health.State(); actual != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}

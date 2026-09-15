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

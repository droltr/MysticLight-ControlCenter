package service

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/droltr/MysticLight-ControlCenter/internal/domain"
)

func TestProfileServiceActivatesAndLoadsProfile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "active.json")
	providers := map[string]struct{}{"openrgb": {}}
	fallback := domain.Profile{
		Name:      "safe",
		Ownership: []domain.Ownership{{Resource: domain.ResourceRGB, Provider: "openrgb"}},
	}
	service, err := NewProfileService(providers, fallback, path)
	if err != nil {
		t.Fatalf("failed to create profile service: %v", err)
	}
	if err := service.Activate(context.Background(), fallback); err != nil {
		t.Fatalf("failed to activate profile: %v", err)
	}
	if err := service.Load(context.Background()); err != nil {
		t.Fatalf("failed to load profile: %v", err)
	}
	if service.Active().Name != "safe" {
		t.Fatalf("unexpected active profile: %q", service.Active().Name)
	}
}

package provider

import (
	"context"
	"testing"

	"github.com/droltr/MysticLight-ControlCenter/internal/domain"
)

func TestRegistryExposesMockProvider(t *testing.T) {
	mock := MockAdapter{
		ProviderName: "openrgb",
		Supported: []domain.Capability{{
			Resource:   domain.ResourceRGB,
			Operations: []domain.Operation{domain.OperationObserve},
		}},
		State: map[string]string{"status": "ready"},
	}
	registry := NewRegistry(mock)

	names := registry.Names()
	if len(names) != 1 || names[0] != "openrgb" {
		t.Fatalf("unexpected provider names: %#v", names)
	}
	state, err := mock.Observe(context.Background())
	if err != nil {
		t.Fatalf("mock observation failed: %v", err)
	}
	if state.(map[string]string)["status"] != "ready" {
		t.Fatalf("unexpected mock state: %#v", state)
	}
}

package service

import (
	"context"
	"testing"
	"time"

	"github.com/droltr/MysticLight-ControlCenter/internal/domain"
	"github.com/droltr/MysticLight-ControlCenter/internal/events"
	"github.com/droltr/MysticLight-ControlCenter/internal/provider"
)

func TestMonitorEmitsSanitizedObservationEvent(t *testing.T) {
	registry := provider.NewRegistry(provider.MockAdapter{
		ProviderName: "openrgb",
		HealthState: domain.ProviderHealth{
			Available: true,
			Error:     "serial=private-device-id",
		},
	})

	eventsFound := NewMonitor(registry, time.Second).Observe(context.Background())
	if len(eventsFound) != 1 {
		t.Fatalf("expected one event, got %d", len(eventsFound))
	}
	if eventsFound[0].Kind != events.ObservationReceived {
		t.Fatalf("unexpected event kind: %s", eventsFound[0].Kind)
	}
	if eventsFound[0].Health.Error != "provider unavailable" {
		t.Fatalf("health error was not sanitized: %q", eventsFound[0].Health.Error)
	}
}

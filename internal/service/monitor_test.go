package service

import (
	"context"
	"errors"
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

func TestMonitorReportsProviderFailureWithoutLeakingError(t *testing.T) {
	registry := provider.NewRegistry(provider.MockAdapter{
		ProviderName: "lcd",
		HealthState:  domain.ProviderHealth{Available: true},
		ObserveError: errors.New("serial=private-device-id"),
	})

	observations := NewMonitor(registry, time.Second).Observe(context.Background())
	if len(observations) != 1 {
		t.Fatalf("expected one event, got %d", len(observations))
	}
	if observations[0].Kind != events.ProviderHealthChanged {
		t.Fatalf("unexpected event kind: %s", observations[0].Kind)
	}
	if observations[0].Health.Available {
		t.Fatal("failed provider must not remain available")
	}
	if observations[0].Health.Error != "provider unavailable" {
		t.Fatalf("provider error was not redacted: %q", observations[0].Health.Error)
	}
}

func TestMonitorReportsTimeoutAsProviderFailure(t *testing.T) {
	registry := provider.NewRegistry(provider.MockAdapter{
		ProviderName: "coolercontrol",
		HealthState:  domain.ProviderHealth{Available: true},
		ObserveDelay: 50 * time.Millisecond,
	})

	observations := NewMonitor(registry, time.Millisecond).Observe(context.Background())
	if observations[0].Kind != events.ProviderHealthChanged {
		t.Fatalf("timeout must produce health change: %s", observations[0].Kind)
	}
	if observations[0].Health.Available {
		t.Fatal("timed-out provider must not remain available")
	}
}

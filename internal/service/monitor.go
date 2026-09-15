package service

import (
	"context"
	"time"

	"github.com/droltr/MysticLight-ControlCenter/internal/domain"
	"github.com/droltr/MysticLight-ControlCenter/internal/events"
	"github.com/droltr/MysticLight-ControlCenter/internal/provider"
)

type Monitor struct {
	registry *provider.Registry
	timeout  time.Duration
}

func NewMonitor(registry *provider.Registry, timeout time.Duration) *Monitor {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &Monitor{registry: registry, timeout: timeout}
}

func (m *Monitor) Observe(ctx context.Context) []events.Event {
	return m.observe(ctx, nil)
}

func (m *Monitor) ObserveAndPublish(ctx context.Context, bus *events.Bus) []events.Event {
	return m.observe(ctx, bus)
}

func (m *Monitor) observe(ctx context.Context, bus *events.Bus) []events.Event {
	observations := make([]events.Event, 0)
	for _, adapter := range m.registry.Adapters() {
		observationCtx, cancel := context.WithTimeout(ctx, m.timeout)
		_, err := adapter.Observe(observationCtx)
		cancel()

		health := adapter.Health(ctx)
		health.Provider = adapter.Name()
		health.Available = err == nil && health.Available
		if err != nil {
			health.Error = "provider observation failed"
		}
		kind := events.ProviderHealthChanged
		if err == nil {
			kind = events.ObservationReceived
		}
		event := events.Event{
			Kind:       kind,
			OccurredAt: time.Now().UTC(),
			Provider:   adapter.Name(),
			Health:     sanitizeHealth(health),
		}
		observations = append(observations, event)
		if bus != nil {
			bus.Publish(event)
		}
	}
	return observations
}

func sanitizeHealth(health domain.ProviderHealth) domain.ProviderHealth {
	if health.Error != "" {
		health.Error = "provider unavailable"
	}
	return health
}

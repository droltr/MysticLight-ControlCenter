package provider

import (
	"context"
	"time"

	"github.com/droltr/MysticLight-ControlCenter/internal/domain"
)

type MockAdapter struct {
	ProviderName string
	Supported    []domain.Capability
	State        any
	HealthState  domain.ProviderHealth
	ObserveError error
	ObserveDelay time.Duration
}

func (m MockAdapter) Name() string {
	return m.ProviderName
}

func (m MockAdapter) Capabilities() []domain.Capability {
	return append([]domain.Capability(nil), m.Supported...)
}

func (m MockAdapter) Health(context.Context) domain.ProviderHealth {
	return m.HealthState
}

func (m MockAdapter) Observe(ctx context.Context) (any, error) {
	if m.ObserveError != nil {
		return nil, m.ObserveError
	}
	if m.ObserveDelay > 0 {
		timer := time.NewTimer(m.ObserveDelay)
		defer timer.Stop()
		select {
		case <-timer.C:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return m.State, nil
}

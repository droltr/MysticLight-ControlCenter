package provider

import (
	"context"

	"github.com/droltr/MysticLight-ControlCenter/internal/domain"
)

type MockAdapter struct {
	ProviderName string
	Supported    []domain.Capability
	State        any
	HealthState  domain.ProviderHealth
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

func (m MockAdapter) Observe(context.Context) (any, error) {
	return m.State, nil
}

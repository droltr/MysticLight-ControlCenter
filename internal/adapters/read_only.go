package adapters

import (
	"context"
	"errors"

	"github.com/droltr/MysticLight-ControlCenter/internal/domain"
	"github.com/droltr/MysticLight-ControlCenter/internal/provider"
)

const (
	OpenRGB       = "openrgb"
	CoolerControl = "coolercontrol"
	LCD           = "lcd"
	Sensors       = "sensors"
	GameFocus     = "game-focus"
)

var ErrNotConfigured = errors.New("provider adapter is not configured")

type NotConfigured struct {
	ProviderName string
	Resource     domain.Resource
}

var _ provider.Adapter = NotConfigured{}

func (a NotConfigured) Name() string {
	return a.ProviderName
}

func (a NotConfigured) Capabilities() []domain.Capability {
	return []domain.Capability{{
		Resource:   a.Resource,
		Operations: []domain.Operation{domain.OperationObserve},
	}}
}

func (a NotConfigured) Health(context.Context) domain.ProviderHealth {
	return domain.ProviderHealth{
		Provider:  a.ProviderName,
		Available: false,
		Error:     "provider unavailable",
	}
}

func (a NotConfigured) Observe(context.Context) (any, error) {
	return nil, ErrNotConfigured
}

func DefaultReadOnlyAdapters() []provider.Adapter {
	return []provider.Adapter{
		NotConfigured{ProviderName: OpenRGB, Resource: domain.ResourceRGB},
		NotConfigured{ProviderName: CoolerControl, Resource: domain.ResourcePWM},
		NotConfigured{ProviderName: LCD, Resource: domain.ResourceLCD},
		NotConfigured{ProviderName: Sensors, Resource: domain.ResourceTelemetry},
		NotConfigured{ProviderName: GameFocus, Resource: domain.ResourceFocus},
	}
}

package provider

import (
	"context"

	"github.com/droltr/MysticLight-ControlCenter/internal/domain"
)

type Adapter interface {
	Name() string
	Capabilities() []domain.Capability
	Health(context.Context) domain.ProviderHealth
	Observe(context.Context) (any, error)
}

type Registry struct {
	providers map[string]Adapter
}

func NewRegistry(adapters ...Adapter) *Registry {
	providers := make(map[string]Adapter, len(adapters))
	for _, adapter := range adapters {
		providers[adapter.Name()] = adapter
	}
	return &Registry{providers: providers}
}

func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}

func (r *Registry) Adapters() []Adapter {
	adapters := make([]Adapter, 0, len(r.providers))
	for _, adapter := range r.providers {
		adapters = append(adapters, adapter)
	}
	return adapters
}

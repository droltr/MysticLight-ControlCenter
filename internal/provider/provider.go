package provider

import (
	"sort"

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
	sort.Strings(names)
	return names
}

func (r *Registry) Adapters() []Adapter {
	adapters := make([]Adapter, 0, len(r.providers))
	for _, adapter := range r.providers {
		adapters = append(adapters, adapter)
	}
	sort.Slice(adapters, func(i, j int) bool {
		return adapters[i].Name() < adapters[j].Name()
	})
	return adapters
}

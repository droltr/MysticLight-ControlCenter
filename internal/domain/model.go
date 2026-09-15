package domain

import (
	"fmt"
	"strings"
	"time"
)

type Resource string

const (
	ResourceRGB       Resource = "rgb"
	ResourcePWM       Resource = "pwm"
	ResourceLCD       Resource = "lcd"
	ResourceTelemetry Resource = "telemetry"
	ResourceFocus     Resource = "focus"
)

type Operation string

const (
	OperationObserve Operation = "observe"
	OperationApply   Operation = "apply"
)

type Capability struct {
	Resource   Resource
	Operations []Operation
}

func (c Capability) Supports(operation Operation) bool {
	for _, supported := range c.Operations {
		if supported == operation {
			return true
		}
	}
	return false
}

type Ownership struct {
	Resource Resource
	Provider string
}

type ProviderHealth struct {
	Provider  string
	Available bool
	LastSeen  time.Time
	Error     string
}

type HealthState string

const (
	HealthAvailable   HealthState = "available"
	HealthDegraded    HealthState = "degraded"
	HealthUnavailable HealthState = "unavailable"
)

func (h ProviderHealth) State() HealthState {
	if !h.Available {
		return HealthUnavailable
	}
	if h.Error != "" {
		return HealthDegraded
	}
	return HealthAvailable
}

type Profile struct {
	Name      string
	Ownership []Ownership
}

func (p Profile) Validate(knownProviders map[string]struct{}) error {
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("profile name is required")
	}
	seen := make(map[Resource]string, len(p.Ownership))
	for _, ownership := range p.Ownership {
		if ownership.Resource == "" {
			return fmt.Errorf("profile %q contains an empty resource", p.Name)
		}
		provider := strings.TrimSpace(ownership.Provider)
		if provider == "" {
			return fmt.Errorf("resource %q has no provider", ownership.Resource)
		}
		if _, ok := knownProviders[provider]; !ok {
			return fmt.Errorf("resource %q references unknown provider %q", ownership.Resource, provider)
		}
		if previous, ok := seen[ownership.Resource]; ok {
			return fmt.Errorf("resource %q has conflicting owners %q and %q", ownership.Resource, previous, provider)
		}
		seen[ownership.Resource] = provider
	}
	return nil
}

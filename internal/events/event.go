package events

import (
	"time"

	"github.com/droltr/MysticLight-ControlCenter/internal/domain"
)

type Kind string

const (
	ProviderHealthChanged Kind = "provider_health_changed"
	ObservationReceived   Kind = "observation_received"
)

type Event struct {
	Kind       Kind
	OccurredAt time.Time
	Provider   string
	Health     domain.ProviderHealth
}

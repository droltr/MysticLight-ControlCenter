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

type Bus struct {
	subscribers []chan Event
	closed      bool
}

func NewBus() *Bus {
	return &Bus{}
}

func (b *Bus) Subscribe(buffer int) <-chan Event {
	if buffer < 1 {
		buffer = 1
	}
	channel := make(chan Event, buffer)
	b.subscribers = append(b.subscribers, channel)
	return channel
}

func (b *Bus) Publish(event Event) {
	if b.closed {
		return
	}
	for _, subscriber := range b.subscribers {
		subscriber <- event
	}
}

func (b *Bus) Close() {
	if b.closed {
		return
	}
	b.closed = true
	for _, subscriber := range b.subscribers {
		close(subscriber)
	}
	b.subscribers = nil
}

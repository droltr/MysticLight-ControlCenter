package events

import "testing"

func TestBusPublishesToSubscriber(t *testing.T) {
	bus := NewBus()
	subscriber := bus.Subscribe(1)
	expected := Event{Kind: ObservationReceived, Provider: "openrgb"}

	bus.Publish(expected)
	actual := <-subscriber
	if actual != expected {
		t.Fatalf("unexpected event: %#v", actual)
	}
}

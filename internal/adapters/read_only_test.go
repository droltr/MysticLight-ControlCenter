package adapters

import (
	"context"
	"testing"

	"github.com/droltr/MysticLight-ControlCenter/internal/domain"
)

func TestDefaultAdaptersAreReadOnlyAndUnavailable(t *testing.T) {
	adapters := DefaultReadOnlyAdapters()
	if len(adapters) != 5 {
		t.Fatalf("expected five default adapters, got %d", len(adapters))
	}
	for _, adapter := range adapters {
		capabilities := adapter.Capabilities()
		if len(capabilities) != 1 || !capabilities[0].Supports(domain.OperationObserve) {
			t.Fatalf("adapter %q is not read-only: %#v", adapter.Name(), capabilities)
		}
		if capabilities[0].Supports(domain.OperationApply) {
			t.Fatalf("adapter %q exposes a write capability", adapter.Name())
		}
		if _, err := adapter.Observe(context.Background()); err != ErrNotConfigured {
			t.Fatalf("adapter %q returned unexpected error: %v", adapter.Name(), err)
		}
	}
}

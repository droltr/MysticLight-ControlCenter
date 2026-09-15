package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/droltr/MysticLight-ControlCenter/internal/adapters"
	"github.com/droltr/MysticLight-ControlCenter/internal/events"
	"github.com/droltr/MysticLight-ControlCenter/internal/provider"
	"github.com/droltr/MysticLight-ControlCenter/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	registry := provider.NewRegistry(adapters.FromEnvironment()...)
	monitor := service.NewMonitor(registry, 5*time.Second)
	bus := events.NewBus()
	_ = bus.Subscribe(len(registry.Names()))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	for _, event := range monitor.ObserveAndPublish(ctx, bus) {
		logger.Info("provider observation", "provider", event.Provider, "kind", event.Kind, "state", event.Health.State())
	}
	<-ctx.Done()
	logger.Info("daemon stopped")
}

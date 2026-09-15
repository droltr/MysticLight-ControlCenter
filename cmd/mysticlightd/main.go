package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/droltr/MysticLight-ControlCenter/internal/adapters"
	"github.com/droltr/MysticLight-ControlCenter/internal/api"
	"github.com/droltr/MysticLight-ControlCenter/internal/domain"
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
	profilePath := os.Getenv("MYSTICLIGHT_PROFILE_PATH")
	if profilePath == "" {
		profilePath = "profile.json"
	}
	profileService, err := service.NewProfileService(providerNames(registry), domain.Profile{Name: "safe"}, profilePath)
	if err != nil {
		logger.Error("create profile service", "error", err)
		return
	}
	apiAddress := os.Getenv("MYSTICLIGHT_API_ADDRESS")
	if apiAddress == "" {
		apiAddress = "127.0.0.1:11888"
	}
	apiServer := &http.Server{Addr: apiAddress, Handler: (api.Server{Registry: registry, ProfileService: profileService}).Handler()}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go monitor.Run(ctx, 30*time.Second, bus)
	go func() {
		if err := apiServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("API server stopped", "error", err)
			stop()
		}
	}()
	<-ctx.Done()
	shutdownContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = apiServer.Shutdown(shutdownContext)
	bus.Close()
	logger.Info("daemon stopped")
}

func providerNames(registry *provider.Registry) map[string]struct{} {
	providers := make(map[string]struct{})
	for _, name := range registry.Names() {
		providers[name] = struct{}{}
	}
	return providers
}

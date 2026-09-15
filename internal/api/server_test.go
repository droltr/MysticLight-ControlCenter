package api

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/droltr/MysticLight-ControlCenter/internal/adapters"
	"github.com/droltr/MysticLight-ControlCenter/internal/domain"
	"github.com/droltr/MysticLight-ControlCenter/internal/provider"
	"github.com/droltr/MysticLight-ControlCenter/internal/service"
)

func TestAPIHealthAndProvidersAreReadOnly(t *testing.T) {
	registry := provider.NewRegistry(adapters.NewOpenRGB())
	profileService, err := service.NewProfileService(
		map[string]struct{}{adapters.OpenRGB: {}},
		domain.Profile{Name: "safe", Ownership: []domain.Ownership{{Resource: domain.ResourceRGB, Provider: adapters.OpenRGB}}},
		filepath.Join(t.TempDir(), "profile.json"),
	)
	if err != nil {
		t.Fatalf("profile service: %v", err)
	}
	handler := Server{Registry: registry, ProfileService: profileService}.Handler()

	for _, path := range []string{"/health", "/providers"} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusOK {
			t.Fatalf("%s returned %d", path, recorder.Code)
		}
	}
}

func TestAPIValidatesAndActivatesProfile(t *testing.T) {
	registry := provider.NewRegistry(adapters.NewOpenRGB())
	profileService, err := service.NewProfileService(
		map[string]struct{}{adapters.OpenRGB: {}},
		domain.Profile{Name: "safe", Ownership: []domain.Ownership{{Resource: domain.ResourceRGB, Provider: adapters.OpenRGB}}},
		filepath.Join(t.TempDir(), "profile.json"),
	)
	if err != nil {
		t.Fatalf("profile service: %v", err)
	}
	handler := Server{Registry: registry, ProfileService: profileService}.Handler()
	body := `{"name":"quiet","ownership":[{"resource":"rgb","provider":"openrgb"}]}`
	request := httptest.NewRequest(http.MethodPost, "/profiles/activate", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("activate returned %d: %s", recorder.Code, recorder.Body.String())
	}
}

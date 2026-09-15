package api

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/droltr/MysticLight-ControlCenter/internal/profile"
	"github.com/droltr/MysticLight-ControlCenter/internal/provider"
	"github.com/droltr/MysticLight-ControlCenter/internal/service"
)

//go:embed web/index.html
var dashboard []byte

type Server struct {
	Registry       *provider.Registry
	ProfileService *service.ProfileService
}

func (s Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.dashboard)
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /providers", s.providers)
	mux.HandleFunc("GET /profiles", s.currentProfile)
	mux.HandleFunc("POST /profiles/validate", s.validateProfile)
	mux.HandleFunc("POST /profiles/activate", s.activateProfile)
	return mux
}

func (s Server) dashboard(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(dashboard)
}

func (s Server) health(writer http.ResponseWriter, _ *http.Request) {
	writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
}

func (s Server) providers(writer http.ResponseWriter, _ *http.Request) {
	items := make([]map[string]any, 0)
	for _, adapter := range s.Registry.Adapters() {
		items = append(items, map[string]any{
			"name":         adapter.Name(),
			"capabilities": adapter.Capabilities(),
		})
	}
	writeJSON(writer, http.StatusOK, items)
}

func (s Server) currentProfile(writer http.ResponseWriter, _ *http.Request) {
	writeJSON(writer, http.StatusOK, s.ProfileService.Active())
}

func (s Server) validateProfile(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	data, err := readBody(request)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err)
		return
	}
	if _, err := profile.ParseJSON(data, providerNames(s.Registry)); err != nil {
		writeError(writer, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(writer, http.StatusOK, map[string]string{"status": "valid"})
}

func (s Server) activateProfile(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	data, err := readBody(request)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err)
		return
	}
	profileValue, err := profile.ParseJSON(data, providerNames(s.Registry))
	if err != nil {
		writeError(writer, http.StatusUnprocessableEntity, err)
		return
	}
	if err := s.ProfileService.Activate(request.Context(), profileValue); err != nil {
		writeError(writer, http.StatusConflict, err)
		return
	}
	writeJSON(writer, http.StatusOK, map[string]string{"status": "activated"})
}

func providerNames(registry *provider.Registry) map[string]struct{} {
	names := make(map[string]struct{})
	for _, name := range registry.Names() {
		names[name] = struct{}{}
	}
	return names
}

func readBody(request *http.Request) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(request.Body, 1<<20+1))
	if err != nil {
		return nil, err
	}
	if len(data) > 1<<20 {
		return nil, fmt.Errorf("request body too large")
	}
	if !json.Valid(data) {
		return nil, fmt.Errorf("request body is not valid JSON")
	}
	return data, nil
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}

func writeError(writer http.ResponseWriter, status int, err error) {
	writeJSON(writer, status, map[string]string{"error": err.Error()})
}

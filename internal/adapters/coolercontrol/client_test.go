package coolercontrol

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientObservesJSONWithBearerToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet || request.URL.Path != "/read-only" {
			t.Fatalf("unexpected request: %s %s", request.Method, request.URL.Path)
		}
		if request.Header.Get("Authorization") != "Bearer test-token" {
			t.Fatalf("authorization header missing")
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"fans":1}`))
	}))
	defer server.Close()

	result, err := New(server.URL, "/read-only", "test-token", time.Second).Observe(context.Background())
	if err != nil {
		t.Fatalf("observe failed: %v", err)
	}
	if string(result.(Observation).Payload) != `{"fans":1}` {
		t.Fatalf("unexpected payload: %s", result.(Observation).Payload)
	}
}

func TestClientRejectsHTTPErrorAndInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/error" {
			http.Error(writer, "failure", http.StatusServiceUnavailable)
			return
		}
		_, _ = writer.Write([]byte("not-json"))
	}))
	defer server.Close()

	if _, err := New(server.URL, "/error", "", time.Second).Observe(context.Background()); err == nil {
		t.Fatal("expected HTTP error")
	}
	if _, err := New(server.URL, "/invalid", "", time.Second).Observe(context.Background()); err == nil {
		t.Fatal("expected invalid JSON error")
	}
}

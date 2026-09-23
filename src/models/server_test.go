package models

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewSimpleServer_StoresAddress(t *testing.T) {
	serverURL := "https://example.com"

	upstream := NewSimpleServer(serverURL)

	if got := upstream.Address(); got != serverURL {
		t.Fatalf("expected address %q, got %q", serverURL, got)
	}
}

func TestSimpleServer_IsAlive_DefaultsFalse(t *testing.T) {
	upstream := NewSimpleServer("https://example.com")

	if upstream.IsAlive() {
		t.Fatal("expected newly created server to start unhealthy until checked")
	}
}

func TestSimpleServer_RefreshHealth_HealthyStatus(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	upstream := NewSimpleServer(backend.URL)
	upstream.RefreshHealth()

	if !upstream.IsAlive() {
		t.Fatal("expected server to be marked healthy after a 200 OK response")
	}
}

func TestSimpleServer_RefreshHealth_UnhealthyStatus(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer backend.Close()

	upstream := NewSimpleServer(backend.URL)
	upstream.RefreshHealth()

	if upstream.IsAlive() {
		t.Fatal("expected server to be marked unhealthy after a 503 response")
	}
}

func TestSimpleServer_Serve_ForwardsRequest(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/hello" {
			t.Fatalf("expected request path /hello, got %s", r.URL.Path)
		}
		rw.WriteHeader(http.StatusAccepted)
		_, _ = rw.Write([]byte("ok"))
	}))
	defer backend.Close()

	upstream := NewSimpleServer(backend.URL)
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rrw := httptest.NewRecorder()

	upstream.Serve(rrw, req)

	if rrw.Code != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d", http.StatusAccepted, rrw.Code)
	}
}

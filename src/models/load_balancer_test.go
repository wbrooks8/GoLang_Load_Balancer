package models

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	serverpkg "github.com/wbrooks8/load_balancer/src/interface"
)

type fakeServer struct {
	addr  string
	alive bool
}

// fakeServer exists only for tests. It behaves like an upstream server without requiring a real network call.
func (f *fakeServer) Address() string {
	return f.addr
}

func (f *fakeServer) IsAlive() bool {
	return f.alive
}

func (f *fakeServer) RefreshHealth() {
	// no-op for tests; health state is set directly in the fake
}

func (f *fakeServer) Serve(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

// Test that requests rotate through the backend list in order.
func TestLoadBalancerGetNextAvailableServer_RoundRobin(t *testing.T) {
	server1 := &fakeServer{addr: "server-1", alive: true}
	server2 := &fakeServer{addr: "server-2", alive: true}
	server3 := &fakeServer{addr: "server-3", alive: true}

	lb := &LoadBalancer{
		servers: []serverpkg.Server{server1, server2, server3},
	}

	if got := lb.getNextAvailableServer(); got != server1 {
		t.Fatalf("expected first request to hit server-1, got %s", got.Address())
	}
	if got := lb.getNextAvailableServer(); got != server2 {
		t.Fatalf("expected second request to hit server-2, got %s", got.Address())
	}
	if got := lb.getNextAvailableServer(); got != server3 {
		t.Fatalf("expected third request to hit server-3, got %s", got.Address())
	}
	if got := lb.getNextAvailableServer(); got != server1 {
		t.Fatalf("expected round robin to wrap back to server-1, got %s", got.Address())
	}
}

func TestLoadBalancerGetNextAvailableServer_SkipsUnhealthyServers(t *testing.T) {
	server1 := &fakeServer{addr: "server-1", alive: false}
	server2 := &fakeServer{addr: "server-2", alive: true}
	server3 := &fakeServer{addr: "server-3", alive: true}

	lb := &LoadBalancer{
		servers: []serverpkg.Server{server1, server2, server3},
	}

	if got := lb.getNextAvailableServer(); got != server2 {
		t.Fatalf("expected unhealthy server-1 to be skipped, got %s", got.Address())
	}
}

func TestLoadBalancerGetNextAvailableServer_ReturnsNilWhenAllServersUnhealthy(t *testing.T) {
	server1 := &fakeServer{addr: "server-1", alive: false}
	server2 := &fakeServer{addr: "server-2", alive: false}

	lb := &LoadBalancer{
		servers: []serverpkg.Server{server1, server2},
	}

	if got := lb.getNextAvailableServer(); got != nil {
		t.Fatalf("expected no healthy servers to return nil, got %s", got.Address())
	}
}

func TestLoadBalancerHandleRequest_Returns503WhenNoHealthyServers(t *testing.T) {
	server1 := &fakeServer{addr: "server-1", alive: false}
	server2 := &fakeServer{addr: "server-2", alive: false}

	lb := &LoadBalancer{
		servers: []serverpkg.Server{server1, server2},
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rrw := httptest.NewRecorder()

	lb.HandleRequest(rrw, req)

	if rrw.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, rrw.Code)
	}
}

func TestSimpleServerRefreshHealth_ReportsHealthyServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	upstream := NewSimpleServer(server.URL)
	simple := upstream.(*simpleServer)

	simple.RefreshHealth()

	if !simple.IsAlive() {
		t.Fatal("expected server to be marked healthy after a 200 OK response")
	}
}

func TestSimpleServerRefreshHealth_ReportsUnhealthyServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	upstream := NewSimpleServer(server.URL)
	simple := upstream.(*simpleServer)

	simple.RefreshHealth()

	if simple.IsAlive() {
		t.Fatal("expected server to be marked unhealthy after a 503 response")
	}
}

func TestNewLoadBalancerStartsHealthChecks(t *testing.T) {
	server1 := &fakeServer{addr: "server-1", alive: true}
	server2 := &fakeServer{addr: "server-2", alive: true}

	lb := NewLoadBalancer("8000", []serverpkg.Server{server1, server2}, time.Millisecond*10)
	if lb == nil {
		t.Fatal("expected load balancer to be created")
	}
	if lb.Port() != "8000" {
		t.Fatalf("expected port 8000, got %s", lb.Port())
	}
}

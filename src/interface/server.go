package server

import "net/http"

type Server interface {
	// Return a stable identifier for the backend.
	Address() string
	// Report whether the backend is currently considered healthy.
	IsAlive() bool
	// Refresh the backend state based on a real health check.
	RefreshHealth()
	// Forward an HTTP request to this backend.
	Serve(rw http.ResponseWriter, r *http.Request)
}

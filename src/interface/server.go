package server

import "net/http"

type Server interface {
	Address() string
	IsAlive() bool
	RefreshHealth()
	Serve(rw http.ResponseWriter, r *http.Request)
}

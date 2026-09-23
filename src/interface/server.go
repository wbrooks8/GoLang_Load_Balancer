package server

import "net/http"

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
	RefreshHealth() 
}

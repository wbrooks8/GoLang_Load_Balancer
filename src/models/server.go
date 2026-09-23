package models

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"

	serverpkg "github.com/wbrooks8/load_balancer/src/interface"
)

type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
	alive bool
	mu    sync.Mutex
}

func NewSimpleServer(addr string) serverpkg.Server {
	serverURL, err := url.Parse(addr)
	handleErr(err)

	return &simpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverURL),
	}
}

func handleErr(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

func (s *simpleServer) Address() string {
	return s.addr
}

func (s *simpleServer) IsAlive() bool {
	// Read the cached health state with locking so concurrent requests do not race.
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.alive
}

func (s *simpleServer) RefreshHealth() {
	// Probe the upstream with a HEAD request and update the cached health flag.
	resp, err := http.Head(s.addr)
	alive := false
	if err == nil {
		defer resp.Body.Close()
		alive = resp.StatusCode >= 200 && resp.StatusCode < 300
	}

	s.mu.Lock()
	s.alive = alive
	s.mu.Unlock()
}

func (s *simpleServer) Serve(rw http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(rw, r)
}

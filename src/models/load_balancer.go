package models

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	serverpkg "github.com/wbrooks8/load_balancer/src/interface"
)

type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []serverpkg.Server
	mu              sync.Mutex
}

func NewLoadBalancer(port string, servers []serverpkg.Server, checkInterval time.Duration) *LoadBalancer {
	lb := &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}

	lb.startHealthChecks(checkInterval)

	return lb
}

func (lb *LoadBalancer) Port() string {
	return lb.port
}

func (lb *LoadBalancer) getNextAvailableServer() serverpkg.Server {
	// Protect the shared round-robin counter so multiple requests do not race.
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// Walk the server list in order and skip any server that is marked unhealthy.
	for i := 0; i < len(lb.servers); i++ {
		// Round-robin index: 0, 1, 2, 0, 1, 2 ... across the available servers.
		server := lb.servers[lb.roundRobinCount%len(lb.servers)]
		lb.roundRobinCount++
		if server.IsAlive() {
			return server
		}
	}
	return nil
}

func (lb *LoadBalancer) HandleRequest(rw http.ResponseWriter, r *http.Request) {
	// Choose a healthy backend before forwarding the incoming request.
	targetServer := lb.getNextAvailableServer()
	if targetServer == nil {
		http.Error(rw, "no healthy servers available", http.StatusServiceUnavailable)
		return
	}

	fmt.Printf("forwarding request to address %q\n", targetServer.Address())
	targetServer.Serve(rw, r)
}

func (lb *LoadBalancer) startHealthChecks(interval time.Duration) {
	// Run one health-check pass immediately so the startup state is not stale.
	lb.checkAll()

	// Keep checking in the background so unhealthy servers can be detected over time.
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			lb.checkAll()
		}
	}()
}

func (lb *LoadBalancer) checkAll() {
	// Refresh each server's health in parallel so one slow server does not block others.
	var wg sync.WaitGroup
	for _, s := range lb.servers {
		wg.Add(1)
		go func(s serverpkg.Server) {
			defer wg.Done()
			s.RefreshHealth()
		}(s)
	}
	wg.Wait()
}

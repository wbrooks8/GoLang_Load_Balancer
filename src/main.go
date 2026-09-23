package main

import (
	"fmt"
	"net/http"

	serverpkg "github.com/wbrooks8/load_balancer/src/interface"
	modelspkg "github.com/wbrooks8/load_balancer/src/models"
)

func main() {
	// Build the list of upstream servers the load balancer will route to.
	servers := []serverpkg.Server{
		modelspkg.NewSimpleServer("https://www.facebook.com"),
		modelspkg.NewSimpleServer("https://www.bing.com"),
		modelspkg.NewSimpleServer("https://www.duckduckgo.com"),
	}

	// Create a balancer with a 15-second health-check interval.
	lb := modelspkg.NewLoadBalancer("8000", servers, 15)

	// Forward every request to the balancer's routing logic.
	http.HandleFunc("/", lb.HandleRequest)

	fmt.Printf("Serving requests at 'localhost: %s'\n", lb.Port())
	if err := http.ListenAndServe(":"+lb.Port(), nil); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}

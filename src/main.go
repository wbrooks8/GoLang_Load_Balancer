package main

import (
	"fmt"
	"net/http"

	serverpkg "github.com/wbrooks8/load_balancer/src/interface"
	modelspkg "github.com/wbrooks8/load_balancer/src/models"
)

func main() {
	servers := []serverpkg.Server{
		modelspkg.NewSimpleServer("https://www.facebook.com"),
		modelspkg.NewSimpleServer("https://www.bing.com"),
		modelspkg.NewSimpleServer("https://www.duckduckgo.com"),
	}

	lb := modelspkg.NewLoadBalancer("8000", servers, 15)

	http.HandleFunc("/", lb.HandleRequest)

	fmt.Printf("Serving requests at 'localhost: %s'\n", lb.Port())
	if err := http.ListenAndServe(":"+lb.Port(), nil); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}

package apiserver

import (
	"net/http"

	"github.com/Cloud-for-You/ipam-for-k8s/pkg/api_server/handlers"
)

func StartServer() {
	// Create handler for functions
	http.HandleFunc("/api/", handlers.HelloWorld)

	http.ListenAndServe(":10080", nil)
}

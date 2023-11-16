package apiserver

import (
	"fmt"
	"net/http"
)

func StartServer() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello world")
	}
	// Create handler for functions
	http.HandleFunc("/api/", handler)

	http.ListenAndServe(":10080", nil)
}

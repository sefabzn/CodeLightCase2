package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Health check endpoint
	http.HandleFunc("/health", healthHandler)

	// Start server
	port := ":8000"
	fmt.Printf("Server starting on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"status": "ok",
	}

	json.NewEncoder(w).Encode(response)
}

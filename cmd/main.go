package main

import (
	"log"
	"net/http"

	"exchangerate/internal/api"
)

func main() {
	http.HandleFunc("/convert", api.ConvertHandler)

	log.Println("Starting Exchange Rate Service on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

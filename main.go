package main

import (
	"net/http"

	"qr-generator/handlers"
	"qr-generator/services"
)

func main() {
	// Initialize service layer
	qrService := services.NewQRService()

	// Initialize handler with dependency injection
	qrHandler := handlers.NewQRHandler(qrService)

	// Setup routes
	http.HandleFunc("/generate", qrHandler.GenerateQRCode)

	// Start server
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/basecase/orchestrator/internal/handlers"
	"github.com/basecase/orchestrator/internal/middleware"
)

func main() {
	mux := http.NewServeMux()

	// Initialize handlers
	sessionHandler := handlers.NewSessionHandler()
	evaluationHandler := handlers.NewEvaluationHandler()

	// Register routes
	mux.HandleFunc("POST /api/v1/sessions", sessionHandler.CreateSession)
	mux.HandleFunc("POST /api/v1/evaluations/submit", evaluationHandler.SubmitEvaluation)

	// Apply middleware
	handler := middleware.RequireUserId(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Starting Orchestrator on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

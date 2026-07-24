package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/basecase/orchestrator/internal/handlers"
	"github.com/basecase/orchestrator/internal/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, relying on system environment variables")
	}

	mux := http.NewServeMux()

	// Initialize handlers
	sessionHandler := handlers.NewSessionHandler()
	evaluationHandler := handlers.NewEvaluationHandler()

	// Register routes
	mux.HandleFunc("POST /api/v1/sessions", sessionHandler.CreateSession)
	mux.HandleFunc("GET /api/v1/sessions/{id}/state", sessionHandler.GetSessionState)
	mux.HandleFunc("POST /api/v1/evaluations/submit", evaluationHandler.SubmitEvaluation)
	mux.HandleFunc("POST /api/v1/evaluations/stream", evaluationHandler.StreamChat)

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

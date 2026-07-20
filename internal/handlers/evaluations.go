package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/basecase/orchestrator/internal/middleware"
	"github.com/basecase/orchestrator/internal/saga"
)

type EvaluationHandler struct {
	saga *saga.EvaluationSaga
}

func NewEvaluationHandler() *EvaluationHandler {
	return &EvaluationHandler{
		saga: saga.NewEvaluationSaga(),
	}
}

func (h *EvaluationHandler) SubmitEvaluation(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Internal Server Error: Missing User ID in context", http.StatusInternalServerError)
		return
	}

	result, err := h.saga.ExecuteSubmitEvaluation(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(result)
}

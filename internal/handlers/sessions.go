package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/basecase/orchestrator/internal/middleware"
	"github.com/basecase/orchestrator/internal/saga"
)

type SessionHandler struct {
	saga *saga.SessionSaga
}

func NewSessionHandler() *SessionHandler {
	return &SessionHandler{
		saga: saga.NewSessionSaga(),
	}
}

func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Internal Server Error: Missing User ID in context", http.StatusInternalServerError)
		return
	}

	var reqBody struct {
		ProblemId string `json:"problemId"`
		Restart   bool   `json:"restart"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.saga.ExecuteCreateSession(r.Context(), userId, reqBody.ProblemId, reqBody.Restart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (h *SessionHandler) GetSessionState(w http.ResponseWriter, r *http.Request) {
	sessionId := r.PathValue("id")
	if sessionId == "" {
		http.Error(w, "Missing session ID", http.StatusBadRequest)
		return
	}

	result, err := h.saga.ExecuteGetSessionState(r.Context(), sessionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

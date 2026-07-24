package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

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

	var reqBody struct {
		SessionId string `json:"sessionId"`
		Code      string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.saga.ExecuteSubmitEvaluation(r.Context(), userId, reqBody.SessionId, reqBody.Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for k, vv := range resp.Header {
		lowerK := strings.ToLower(k)
		if lowerK == "transfer-encoding" || lowerK == "content-length" {
			continue
		}
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)

	flusher, canFlush := w.(http.Flusher)
	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			w.Write(buf[:n])
			if canFlush {
				flusher.Flush()
			}
		}
		if err != nil {
			break
		}
	}
}

func (h *EvaluationHandler) StreamChat(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		ThreadId string                 `json:"threadId"`
		Input    map[string]interface{} `json:"input"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.saga.ExecuteStreamChat(r.Context(), reqBody.ThreadId, reqBody.Input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for k, vv := range resp.Header {
		lowerK := strings.ToLower(k)
		if lowerK == "transfer-encoding" || lowerK == "content-length" {
			continue
		}
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)

	flusher, canFlush := w.(http.Flusher)
	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			w.Write(buf[:n])
			if canFlush {
				flusher.Flush()
			}
		}
		if err != nil {
			break
		}
	}
}

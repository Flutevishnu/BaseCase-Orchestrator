package saga

import (
	"context"
	"fmt"
	"log"

	"github.com/basecase/orchestrator/internal/clients"
)

type SessionSaga struct {
	express   *clients.ExpressClient
	langgraph *clients.LangGraphClient
}

func NewSessionSaga() *SessionSaga {
	return &SessionSaga{
		express:   clients.NewExpressClient(),
		langgraph: clients.NewLangGraphClient(),
	}
}

func (s *SessionSaga) ExecuteCreateSession(ctx context.Context, userId string, problemId string, restart bool) (map[string]interface{}, error) {
	sessionRes, err := s.express.CreateSession(ctx, userId, problemId, restart)
	if err != nil {
		return nil, fmt.Errorf("failed to create/resume session in Express: %w", err)
	}

	dataObj, ok := sessionRes["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format from Express")
	}

	sessionId, _ := dataObj["threadId"].(string)
	isNew, _ := dataObj["isNew"].(bool)

	if isNew {
		_, err := s.langgraph.CreateThread(ctx, userId, sessionId)
		if err != nil {
			log.Printf("LangGraph thread creation failed, compensating session %s", sessionId)
			compErr := s.express.MarkSessionFailed(ctx, sessionId)
			if compErr != nil {
				log.Printf("CRITICAL: Compensating transaction failed for session %s: %v", sessionId, compErr)
			}
			return nil, fmt.Errorf("failed to initialize AI thread, session aborted: %w", err)
		}

		// Seed initial memory state for the AI
		problemTags, okTags := dataObj["problemTags"].([]interface{})
		userMemory, okMem := dataObj["userMemory"].(map[string]interface{})
		
		if okTags && okMem {
			values := map[string]interface{}{
				"problem_tags":          problemTags,
				"technical_blindspots":  userMemory["technicalBlindspots"],
				"behavioral_traits":     userMemory["behavioralTraits"],
				"past_related_problems": userMemory["pastRelatedProblems"],
			}
			err = s.langgraph.UpdateState(ctx, sessionId, values)
			if err != nil {
				log.Printf("Warning: Failed to seed LangGraph state for session %s: %v", sessionId, err)
			}
		}
	}

	return sessionRes, nil
}

func (s *SessionSaga) ExecuteGetSessionState(ctx context.Context, sessionId string) (map[string]interface{}, error) {
	return s.langgraph.GetState(ctx, sessionId)
}

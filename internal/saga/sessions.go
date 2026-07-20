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

func (s *SessionSaga) ExecuteCreateSession(ctx context.Context, userId string) (map[string]interface{}, error) {
	sessionRes, err := s.express.CreateSession(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to create/resume session in Express: %w", err)
	}

	sessionId, _ := sessionRes["id"].(string)
	isNew, _ := sessionRes["isNew"].(bool)

	if isNew {
		_, err := s.langgraph.CreateThread(ctx, userId)
		if err != nil {

			log.Printf("LangGraph thread creation failed, compensating session %s", sessionId)
			compErr := s.express.MarkSessionFailed(ctx, sessionId)
			if compErr != nil {
				log.Printf("CRITICAL: Compensating transaction failed for session %s: %v", sessionId, compErr)
			}
			return nil, fmt.Errorf("failed to initialize AI thread, session aborted: %w", err)
		}
	}

	return sessionRes, nil
}

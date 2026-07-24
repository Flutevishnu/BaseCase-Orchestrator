package saga

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/basecase/orchestrator/internal/clients"
)

type EvaluationSaga struct {
	express   *clients.ExpressClient
	langgraph *clients.LangGraphClient
}

func NewEvaluationSaga() *EvaluationSaga {
	return &EvaluationSaga{
		express:   clients.NewExpressClient(),
		langgraph: clients.NewLangGraphClient(),
	}
}

func (s *EvaluationSaga) Express() *clients.ExpressClient {
	return s.express
}

func (s *EvaluationSaga) ExecuteSubmitEvaluation(ctx context.Context, userId string, sessionId string, code string) (*http.Response, error) {

	submissionRes, err := s.express.CreateSubmission(ctx, userId, sessionId, code)
	if err != nil {
		return nil, fmt.Errorf("failed to create submission in Express: %w", err)
	}

	submissionId, _ := submissionRes["id"].(string)

	resp, err := s.langgraph.StreamEvaluationRun(ctx, sessionId, submissionId, code)
	if err != nil {
		log.Printf("LangGraph evaluation start failed, compensating submission %s", submissionId)
		compErr := s.express.MarkSubmissionFailed(ctx, submissionId)
		if compErr != nil {
			log.Printf("CRITICAL: Compensating transaction failed for submission %s: %v", submissionId, compErr)
		}
		return nil, fmt.Errorf("failed to start LangGraph evaluation: %w", err)
	}

	return resp, nil
}

func (s *EvaluationSaga) ExecuteStreamChat(ctx context.Context, threadId string, input map[string]interface{}) (*http.Response, error) {
	return s.langgraph.StreamChat(ctx, threadId, input)
}

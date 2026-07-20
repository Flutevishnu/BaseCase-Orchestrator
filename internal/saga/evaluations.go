package saga

import (
	"context"
	"fmt"
	"log"

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

func (s *EvaluationSaga) ExecuteSubmitEvaluation(ctx context.Context, userId string) (map[string]interface{}, error) {

	submissionRes, err := s.express.CreateSubmission(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to create submission in Express: %w", err)
	}

	submissionId, _ := submissionRes["id"].(string)

	_, err = s.langgraph.StreamEvaluationRun(ctx, userId, submissionId)
	if err != nil {
		log.Printf("LangGraph evaluation start failed, compensating submission %s", submissionId)
		compErr := s.express.MarkSubmissionFailed(ctx, submissionId)
		if compErr != nil {
			log.Printf("CRITICAL: Compensating transaction failed for submission %s: %v", submissionId, compErr)
		}
		return nil, fmt.Errorf("failed to start LangGraph evaluation: %w", err)
	}

	return submissionRes, nil
}

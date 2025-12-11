package writing

import (
	"context"
	"encoding/json"
	"fmt"
	dto "zenthara/internal/domain/dto/writing"
	"zenthara/internal/domain/models"
	"zenthara/internal/services/gpt"
	"zenthara/internal/services/prompts"
	"zenthara/internal/utils/converter"

	"github.com/rs/zerolog"
)

type EvaluationService interface {
	EvaluateAnswer(ctx context.Context, req dto.EvaluationRequest) (json.RawMessage, error)
}

type evaluationService struct {
	gptClient     gpt.Client
	promptManager *prompts.WritingPrompts
	logger        zerolog.Logger
}

func NewEvaluationService(
	gptClient gpt.Client,
	promptManager *prompts.WritingPrompts,
	logger zerolog.Logger,
) EvaluationService {
	return &evaluationService{
		gptClient:     gptClient,
		promptManager: promptManager,
		logger:        logger.With().Str("service", "evaluation").Logger(),
	}
}

func (s *evaluationService) EvaluateAnswer(ctx context.Context, req dto.EvaluationRequest) (json.RawMessage, error) {
	logger := s.logger.With().
		Str("testType", string(req.TestType)).
		RawJSON("data", req.Data).
		Logger()

	logger.Debug().Msg("Starting evaluation")

	genericReq, err := converter.ToEvaluationPromptData(req)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to convert request to evaluation prompt data")
		return nil, fmt.Errorf("failed to convert request to evaluation prompt data: %w", err)
	}

	promptData := models.EvaluationPromptData(genericReq)
	prompt, err := s.promptManager.GetEvaluationPrompt(promptData)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get evaluation prompt")
		return nil, fmt.Errorf("failed to get evaluation prompt: %w", err)
	}

	response, err := s.gptClient.Complete(ctx, prompt)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to evaluate answer with GPT")
		return nil, fmt.Errorf("failed to evaluate answer: %w", err)
	}

	if err := s.validateResponse(response); err != nil {
		logger.Error().Err(err).RawJSON("response", response).Msg("Invalid response structure from GPT")
		return nil, fmt.Errorf("invalid response structure: %w", err)
	}

	return response, nil
}

func (s *evaluationService) validateResponse(response json.RawMessage) error {
	// TODO: Implement response validation
	return nil
}

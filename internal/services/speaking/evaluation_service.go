package speaking

import (
	"context"
	"encoding/json"
	"fmt"
	speakingdto "zenthara/internal/domain/dto/speaking"
	"zenthara/internal/domain/models"
	"zenthara/internal/services/gpt"
	"zenthara/internal/services/prompts"
	"zenthara/internal/utils/converter"

	"github.com/rs/zerolog"
)

type EvaluationService interface {
	EvaluateAnswer(ctx context.Context, req speakingdto.EvaluationRequest) (json.RawMessage, error)
}

type evaluationService struct {
	gptClient     gpt.Client
	promptManager *prompts.SpeakingPrompts
	logger        zerolog.Logger
}

func NewEvaluationService(
	gptClient gpt.Client,
	promptManager *prompts.SpeakingPrompts,
	logger zerolog.Logger,
) EvaluationService {
	return &evaluationService{
		gptClient:     gptClient,
		promptManager: promptManager,
		logger:        logger.With().Str("service", "evaluation").Logger(),
	}
}

func (s *evaluationService) EvaluateAnswer(ctx context.Context, req speakingdto.EvaluationRequest) (json.RawMessage, error) {
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
	var respMap map[string]any
	if err := json.Unmarshal(response, &respMap); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	requiredFields := []string{
		"sectionResults",
		"overallBandScore",
		"generalFeedback",
	}

	for _, field := range requiredFields {
		if _, exists := respMap[field]; !exists {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	sectionResults, ok := respMap["sectionResults"].(map[string]any)
	if !ok {
		return fmt.Errorf("invalid sectionResults structure")
	}

	requiredSections := []string{
		"fluencyAndCoherence",
		"lexicalResource",
		"grammaticalRangeAndAccuracy",
		"pronunciation",
	}

	for _, section := range requiredSections {
		if _, exists := sectionResults[section]; !exists {
			return fmt.Errorf("missing required section: %s", section)
		}
	}

	return nil
}

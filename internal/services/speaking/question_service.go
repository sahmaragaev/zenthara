package speaking

import (
	"context"
	"encoding/json"
	"fmt"
	speakingservices "zenthara/internal/domain/dto/speaking"
	testtype "zenthara/internal/domain/enums/testtype/speaking"
	"zenthara/internal/domain/models"
	"zenthara/internal/services/gpt"
	"zenthara/internal/services/prompts"

	"github.com/rs/zerolog"
)

type QuestionService interface {
	GenerateQuestions(ctx context.Context, req speakingservices.QuestionRequest) (json.RawMessage, error)
}

type questionService struct {
	gptClient     gpt.Client
	promptManager *prompts.SpeakingPrompts
	logger        zerolog.Logger
}

func NewQuestionService(
	gptClient gpt.Client,
	promptManager *prompts.SpeakingPrompts,
	logger zerolog.Logger,
) QuestionService {
	return &questionService{
		gptClient:     gptClient,
		promptManager: promptManager,
		logger:        logger.With().Str("writing_service", "question").Logger(),
	}
}

func (s *questionService) GenerateQuestions(ctx context.Context, req speakingservices.QuestionRequest) (json.RawMessage, error) {
	logger := s.logger.With().
		Str("testType", string(req.TestType)).
		Interface("includeTopics", req.IncludeTopics).
		Interface("excludeTopics", req.ExcludeTopics).
		Logger()

	logger.Debug().Msg("Generating questions")

	prompt, err := s.promptManager.GetPrompt(req.TestType, models.QuestionPromptData{
		IncludeTopics: req.IncludeTopics,
		ExcludeTopics: req.ExcludeTopics,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get prompt: %w", err)
	}

	response, err := s.gptClient.Complete(ctx, prompt)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to generate questions from GPT")
		return nil, fmt.Errorf("failed to generate questions: %w", err)
	}

	if err := s.validateResponse(response, req.TestType); err != nil {
		logger.Error().Err(err).RawJSON("response", response).Msg("Invalid response structure from GPT")
		return nil, fmt.Errorf("invalid response structure: %w", err)
	}

	return response, nil
}

func (s *questionService) validateResponse(response json.RawMessage, testType testtype.TestType) error {
	var respMap map[string]any
	if err := json.Unmarshal(response, &respMap); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	requiredFields := getRequiredFields(testType)
	for _, field := range requiredFields {
		if _, exists := respMap[field]; !exists {
			return fmt.Errorf("missing %s in response", field)
		}
	}

	return nil
}

func getRequiredFields(testType testtype.TestType) []string {
	switch testType {
	case testtype.Part1Only:
		return []string{"part1"}
	case testtype.Part2Only:
		return []string{"part2"}
	case testtype.Part2And3:
		return []string{"part2", "part3"}
	case testtype.FullTest:
		return []string{"part1", "part2", "part3"}
	default:
		return []string{}
	}
}

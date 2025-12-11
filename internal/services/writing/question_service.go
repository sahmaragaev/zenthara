package writing

import (
	"context"
	"encoding/json"
	"fmt"
	writingservices "zenthara/internal/domain/dto/writing"
	testtype "zenthara/internal/domain/enums/testtype/writing"
	"zenthara/internal/domain/models"
	"zenthara/internal/services/gpt"
	"zenthara/internal/services/prompts"

	"github.com/rs/zerolog"
)

type QuestionService interface {
	GenerateQuestions(ctx context.Context, req writingservices.QuestionRequest) (json.RawMessage, error)
}

type questionService struct {
	gptClient     gpt.Client
	promptManager *prompts.WritingPrompts
	logger        zerolog.Logger
}

func NewQuestionService(
	gptClient gpt.Client,
	promptManager *prompts.WritingPrompts,
	logger zerolog.Logger,
) QuestionService {
	return &questionService{
		gptClient:     gptClient,
		promptManager: promptManager,
		logger:        logger.With().Str("speaking_service", "question").Logger(),
	}
}

func (s *questionService) GenerateQuestions(ctx context.Context, req writingservices.QuestionRequest) (json.RawMessage, error) {
	logger := s.logger.With().
		Str("testType", string(req.TestType)).
		Interface("includeTopics", req.IncludeTopics).
		Interface("excludeTopics", req.ExcludeTopics).
		Logger()

	logger.Debug().Msg("Generating question")

	taskType := req.TaskType.String()
	prompt, err := s.promptManager.GetPrompt(req.TestType, models.QuestionPromptData{
		IncludeTopics: req.IncludeTopics,
		ExcludeTopics: req.ExcludeTopics,
		TaskType:      &taskType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get prompt: %w", err)
	}

	response, err := s.gptClient.Complete(ctx, prompt)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to generate question from GPT")
		return nil, fmt.Errorf("failed to generate question: %w", err)
	}

	if err := s.validateResponse(response, req.TestType); err != nil {
		return nil, fmt.Errorf("failed to validate response: %w", err)
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
	case testtype.Task2:
		return []string{"statement", "question"}
	default:
		return []string{}
	}
}

package prompts

import (
	"fmt"
	"zenthara/internal/domain/enums/prompttype"
	testtype "zenthara/internal/domain/enums/testtype/speaking"
	"zenthara/internal/domain/models"
)

type SpeakingPrompts struct {
	store *PromptStore
}

func NewSpeakingPrompts(store *PromptStore) *SpeakingPrompts {
	return &SpeakingPrompts{
		store: store,
	}
}

func (p *SpeakingPrompts) GetPrompt(testType testtype.TestType, data models.QuestionPromptData) (string, error) {
	return p.getCommonPrompt(testType, prompttype.SpeakingQuestions, data)
}

func (p *SpeakingPrompts) GetEvaluationPrompt(data models.EvaluationPromptData) (string, error) {
	testType := testtype.TestType(data.TestType)

	return p.getCommonPrompt(testType, prompttype.SpeakingEvaluation, data)
}

func (p *SpeakingPrompts) getCommonPrompt(testType testtype.TestType, promptType prompttype.PromptType, data any) (string, error) {
	var promptName string
	switch testType {
	case testtype.Part1Only:
		promptName = "part1_only"
	case testtype.Part2Only:
		promptName = "part2_only"
	case testtype.Part2And3:
		promptName = "part2_and_3"
	case testtype.FullTest:
		promptName = "full_test"
	default:
		return "", fmt.Errorf("invalid test type: %s", testType)
	}

	return p.store.GetPrompt(promptType, promptName, data)
}

package prompts

import (
	"fmt"
	"zenthara/internal/domain/enums/prompttype"
	testtype "zenthara/internal/domain/enums/testtype/writing"
	"zenthara/internal/domain/models"
)

type WritingPrompts struct {
	store *PromptStore
}

func NewWritingPrompts(store *PromptStore) *WritingPrompts {
	return &WritingPrompts{
		store: store,
	}
}

func (p *WritingPrompts) GetPrompt(testType testtype.TestType, data models.QuestionPromptData) (string, error) {
	return p.getCommonPrompt(testType, prompttype.WritingQuestions, data)
}

func (p *WritingPrompts) GetEvaluationPrompt(data models.EvaluationPromptData) (string, error) {
	testType := testtype.TestType(data.TestType)

	return p.getCommonPrompt(testType, prompttype.WritingEvaluation, data)
}

func (p *WritingPrompts) getCommonPrompt(testType testtype.TestType, promptType prompttype.PromptType, data any) (string, error) {
	var promptName string
	switch testType {
	case testtype.Task2:
		promptName = "task2"
	default:
		return "", fmt.Errorf("invalid test type: %s", testType)
	}

	return p.store.GetPrompt(promptType, promptName, data)
}

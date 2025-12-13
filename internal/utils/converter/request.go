package converter

import (
	"encoding/json"
	"fmt"
	"zenthara/internal/domain/dto/speaking"
	"zenthara/internal/domain/dto/writing"
	"zenthara/internal/domain/models"
)

func FromSpeakingRequest(req *speaking.EvaluationRequest) models.EvaluationPromptData {
	prompt := models.EvaluationPromptData{
		Category: req.Category,
		TestType: string(req.TestType),
		TaskType: "",
		Data:     req.Data,
	}

	var payload struct {
		Question  string `json:"question"`
		Answer    string `json:"answer"`
		WordCount int    `json:"wordCount"`
	}
	if err := json.Unmarshal(req.Data, &payload); err == nil {
		prompt.Question = payload.Question
		prompt.Answer = payload.Answer
		prompt.WordCount = payload.WordCount
	}

	return prompt
}

func FromWritingRequest(req *writing.EvaluationRequest) models.EvaluationPromptData {
	prompt := models.EvaluationPromptData{
		Category: req.Category,
		TestType: string(req.TestType),
		TaskType: string(req.TaskType),
		Data:     req.Data,
	}

	var payload struct {
		Question  string `json:"question"`
		Answer    string `json:"answer"`
		WordCount int    `json:"wordCount"`
	}
	if err := json.Unmarshal(req.Data, &payload); err == nil {
		prompt.Question = payload.Question
		prompt.Answer = payload.Answer
		prompt.WordCount = payload.WordCount
	}

	return prompt
}

func ToEvaluationPromptData(req any) (models.EvaluationPromptData, error) {
	switch v := req.(type) {
	case *speaking.EvaluationRequest:
		return FromSpeakingRequest(v), nil
	case *writing.EvaluationRequest:
		return FromWritingRequest(v), nil
	case speaking.EvaluationRequest:
		return FromSpeakingRequest(&v), nil
	case writing.EvaluationRequest:
		return FromWritingRequest(&v), nil
	default:
		return models.EvaluationPromptData{}, fmt.Errorf("unsupported request type: %T", req)
	}
}

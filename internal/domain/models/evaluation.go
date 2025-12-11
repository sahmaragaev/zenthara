package models

import "encoding/json"

type TestTypeCategory string

const (
	TestTypeSpeaking TestTypeCategory = "speaking"
	TestTypeWriting  TestTypeCategory = "writing"
)

type EvaluationPromptData struct {
	Category  TestTypeCategory
	TestType  string
	TaskType  string
	Question  string
	Answer    string
	WordCount int
	Data      json.RawMessage
}

package writing

import (
	"encoding/json"
	"errors"
	"fmt"
	writing "zenthara/internal/domain/enums/testtype/writing"
	"zenthara/internal/domain/models"
)

type EvaluationRequest struct {
	Category models.TestTypeCategory `json:"category" binding:"required"`
	TestType writing.TestType        `json:"test_type" binding:"required"`
	TaskType writing.TaskType        `json:"task_type" binding:"required"`
	Data     json.RawMessage         `json:"data" binding:"required"`
}

func (r *EvaluationRequest) Validate() error {
	if r.Category == "" {
		return errors.New("category is required")
	}

	if r.TestType == "" {
		return errors.New("test type is required")
	}

	if len(r.Data) == 0 {
		return errors.New("data is required")
	}

	switch r.Category {
	case models.TestTypeWriting:
		if !writing.TestType(r.TestType).IsValid() {
			return fmt.Errorf("invalid writing test type: %s", r.TestType)
		}
	default:
		return fmt.Errorf("invalid category: %s", r.Category)
	}

	return nil
}

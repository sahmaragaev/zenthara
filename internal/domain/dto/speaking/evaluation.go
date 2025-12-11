package speaking

import (
	"encoding/json"
	"errors"
	"fmt"
	speaking "zenthara/internal/domain/enums/testtype/speaking"
	"zenthara/internal/domain/models"
)

type EvaluationRequest struct {
	Category models.TestTypeCategory `json:"category" binding:"required"`
	TestType speaking.TestType       `json:"test_type" binding:"required"`
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
	case models.TestTypeSpeaking:
		if !speaking.TestType(r.TestType).IsValid() {
			return fmt.Errorf("invalid speaking test type: %s", r.TestType)
		}
	default:
		return fmt.Errorf("invalid category: %s", r.Category)
	}

	return nil
}

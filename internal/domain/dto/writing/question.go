package writing

import "zenthara/internal/domain/enums/testtype/writing"

type (
	QuestionRequest struct {
		TestType      writing.TestType `json:"test_type" binding:"required,oneof=task2"`
		TaskType      writing.TaskType `json:"task_type,omitempty" binding:"omitempty,oneof=opinion advantage_disadvantage problem_solution discuss_both_views"`
		IncludeTopics []string         `json:"include_topics,omitempty"`
		ExcludeTopics []string         `json:"exclude_topics,omitempty"`
	}

	QuestionResponse struct {
		Statement string `json:"statement"`
		Question  string `json:"question"`
	}
)

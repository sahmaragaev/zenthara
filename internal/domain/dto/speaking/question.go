package speaking

import testtype "zenthara/internal/domain/enums/testtype/speaking"

type (
	QuestionRequest struct {
		TestType      testtype.TestType `json:"test_type" binding:"required,oneof=part1_only part2_only part2_and_3 full_test"`
		IncludeTopics []string          `json:"include_topics,omitempty"`
		ExcludeTopics []string          `json:"exclude_topics,omitempty"`
	}

	Part1Question struct {
		Questions []string `json:"questions"`
	}

	Part2Question struct {
		MainQuestion string   `json:"main_question"`
		Cues         []string `json:"cues"`
	}

	Part3Question struct {
		Questions []string `json:"questions"`
	}

	QuestionResponse struct {
		Part1 *Part1Question `json:"part1,omitempty"`
		Part2 *Part2Question `json:"part2,omitempty"`
		Part3 *Part3Question `json:"part3,omitempty"`
	}
)

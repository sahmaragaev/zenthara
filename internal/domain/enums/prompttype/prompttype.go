package prompttype

type PromptType string

const (
	SpeakingQuestions  PromptType = "speaking_questions"
	SpeakingEvaluation PromptType = "speaking_evaluation"
	WritingQuestions   PromptType = "writing_questions"
	WritingEvaluation  PromptType = "writing_evaluation"	
)

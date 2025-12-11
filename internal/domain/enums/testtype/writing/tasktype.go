package writing

type TaskType string

const (
	Opinion TaskType = "opinion"
	AdvantageDisadvantage TaskType = "advantage_disadvantage"
	ProblemSolution TaskType = "problem_solution"
	DiscussBothViews TaskType = "discuss_both_views"
)

func (t TaskType) String() string {
	return string(t)
}

func (t TaskType) IsValid() bool {
	switch t {
	case Opinion, AdvantageDisadvantage, ProblemSolution, DiscussBothViews:
		return true
	default:
		return false
	}
}
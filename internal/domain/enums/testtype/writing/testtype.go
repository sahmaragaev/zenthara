package writing

type TestType string

const (
	Task2 TestType = "task2"
)

func (t TestType) String() string {
	return string(t)
}

func (t TestType) IsValid() bool {
	switch t {
	case Task2:
		return true
	default:
		return false
	}
}

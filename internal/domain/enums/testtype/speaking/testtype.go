package speaking

type TestType string

const (
	Part1Only TestType = "part1_only"
	Part2Only TestType = "part2_only"
	Part2And3 TestType = "part2_and_3"
	FullTest  TestType = "full_test"
)

func (t TestType) String() string {
	return string(t)
}

func (t TestType) IsValid() bool {
	switch t {
	case Part1Only, Part2Only, Part2And3, FullTest:
		return true
	default:
		return false
	}
}

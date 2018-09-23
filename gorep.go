package main

type StateMachine struct {
	acceptingStates  map[int]bool
	stateTransitions map[int]map[rune]int
}

func Compile(reStr string) (sm StateMachine, err error) {
	return *new(StateMachine), nil
}

func (sm StateMachine) Test(str string) bool {
	return false
}

func (sm StateMachine) Match(str string) (ans string, matchFound bool) {
	return "", false
}

func main() {
}

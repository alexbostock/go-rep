package compiler

type StateMachine struct {
	acceptingStates  map[int]bool
	stateTransitions map[int]map[rune]int
}

func (tree ast) BuildStateMachine() (sm StateMachine, err error) {
	return *new(StateMachine), nil
}

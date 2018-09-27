package compiler

import (
	"errors"
)

type StateMachine struct {
	AcceptingStates  map[int]bool
	StateTransitions map[int]map[rune]int
}

func (tree ast) BuildStateMachine() (sm StateMachine, err error) {
	switch tree.label {
	case empty:
		sm = *new(StateMachine)
		sm.AcceptingStates = make(map[int]bool)
		sm.StateTransitions = make(map[int]map[rune]int)

		sm.AcceptingStates[0] = true
		sm.StateTransitions[0] = make(map[rune]int)

		for _, child := range tree.children {
			subMachine, err := child.BuildStateMachine()
			if err != nil {
				return sm, err
			}
			sm.append(subMachine)
		}

		return sm, nil
	case literal:
		sm = *new(StateMachine)
		sm.AcceptingStates = make(map[int]bool)
		sm.StateTransitions = make(map[int]map[rune]int)

		sm.AcceptingStates[1] = true

		sm.StateTransitions[0] = make(map[rune]int)
		sm.StateTransitions[1] = make(map[rune]int)

		sm.StateTransitions[0][tree.literal] = 1

		return sm, nil
	case plus:
		if len(tree.children) != 1 {
			return *new(StateMachine), errors.New("Invalid tree: plus node should have exactly one child")
		}

		sm, err := tree.children[0].BuildStateMachine()
		if err != nil {
			return *new(StateMachine), err
		}

		// Accepting states should act as start states
		for state, accepting := range sm.AcceptingStates {
			if !accepting {
				continue
			}

			for char, nextState := range sm.StateTransitions[0] {
				sm.StateTransitions[state][char] = nextState
			}
		}

		return sm, nil
	case optional:
		if len(tree.children) != 1 {
			return *new(StateMachine), errors.New("Invalid tree: ? node should have exactly one child")
		}

		sm, err := tree.children[0].BuildStateMachine()
		if err != nil {
			return *new(StateMachine), err
		}

		// The start state is an accepting state
		sm.AcceptingStates[0] = true

		return sm, err
	case star:
		if len(tree.children) != 1 {
			return *new(StateMachine), errors.New("Invalid tree: * node should have exactly one child")
		}

		sm, err := tree.children[0].BuildStateMachine()
		if err != nil {
			return sm, err
		}

		// The start state is an accepting state
		sm.AcceptingStates[0] = true

		// All transitions from the start state are also transitions from all accepting states

		for char, nextState := range sm.StateTransitions[0] {
			for state, accepting := range sm.AcceptingStates {
				if !accepting {
					continue
				}

				// TODO: this won't work if state already has a transition for char (to a different state)
				sm.StateTransitions[state][char] = nextState
			}
		}

		return sm, nil
	case anyof:
		if len(tree.children) == 0 {
			return *new(StateMachine), errors.New("Invalid tree: [] node must have at least one child")
		}

		// TODO

		return *new(StateMachine), errors.New("Not yet implemented!")
	default:
		return *new(StateMachine), errors.New("This should never happen! Unrecognised ast node")
	}
}

func (sm StateMachine) append(sn StateMachine) {
	offset := len(sm.AcceptingStates) - 1

	for state, accepting := range sn.AcceptingStates {
		sm.AcceptingStates[state+offset] = accepting
	}

	for state, transitions := range sn.StateTransitions {
		sm.StateTransitions[state+offset] = make(map[rune]int)
		for char, nextState := range transitions {
			sm.StateTransitions[state+offset][char] = nextState + offset
		}
	}

	// This has the same problem described above in case star:

	for i := 0; i <= offset; i++ {
		if sm.AcceptingStates[i] {
			sm.AcceptingStates[i] = false
			for char, nextState := range sn.StateTransitions[0] {
				sm.StateTransitions[i][char] = nextState + offset
			}
		}
	}
}

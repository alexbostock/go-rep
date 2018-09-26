package compiler

import (
	"errors"
)

type StateMachine struct {
	acceptingStates  map[int]bool
	stateTransitions map[int]map[rune]int
}

func (tree ast) BuildStateMachine() (sm StateMachine, err error) {
	switch tree.label {
	case empty:
		sm = *new(StateMachine)
		sm.acceptingStates = make(map[int]bool)
		sm.stateTransitions = make(map[int]map[rune]int)

		sm.acceptingStates[0] = true
		sm.stateTransitions[0] = make(map[rune]int)

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
		sm.acceptingStates = make(map[int]bool)
		sm.stateTransitions = make(map[int]map[rune]int)

		sm.acceptingStates[1] = true

		sm.stateTransitions[0] = make(map[rune]int)
		sm.stateTransitions[1] = make(map[rune]int)

		sm.stateTransitions[0][tree.literal] = 1

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
		for state, accepting := range sm.acceptingStates {
			if !accepting {
				continue
			}

			for char, nextState := range sm.stateTransitions[0] {
				sm.stateTransitions[state][char] = nextState
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
		sm.acceptingStates[0] = true

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
		sm.acceptingStates[0] = true

		// All transitions from the start state are also transitions from all accepting states

		for char, nextState := range sm.stateTransitions[0] {
			for state, accepting := range sm.acceptingStates {
				if !accepting {
					continue
				}

				// TODO: this won't work if state already has a transition for char (to a different state)
				sm.stateTransitions[state][char] = nextState
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
	// TODO
}

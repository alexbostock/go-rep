package main

import (
	"fmt"
	"os"
	"strings"
)

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
	if len(os.Args) != 4 {
		fmt.Println("Usage: gorep (test|match) <regex> <string>")
		return
	}

	switch strings.ToLower(os.Args[1]) {
	case "test":
		if sm, err := Compile(os.Args[2]); err == nil {
			fmt.Println(sm.Test(os.Args[3]))
		} else {
			fmt.Println("Invalid regex:", os.Args[2])
		}
	case "match":
		if sm, err := Compile(os.Args[2]); err == nil {
			if ans, ok := sm.Match(os.Args[3]); ok {
				fmt.Println("Match found:", ans)
			}
			fmt.Println("No match found")
		} else {
			fmt.Println("Invalid regex:", os.Args[2])
		}
	default:
		fmt.Println("Invalid option:", os.Args[1])
		fmt.Println("Usage: gorep (test|match) <regex> <string>")
	}
}

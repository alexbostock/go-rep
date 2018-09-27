package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexbostock/go-rep/compiler"
)

func Compile(reStr string) (sm compiler.StateMachine, err error) {
	ast, err := compiler.Parse(reStr)
	if err != nil {
		return *new(compiler.StateMachine), err
	}

	return ast.BuildStateMachine()
}

func Test(sm compiler.StateMachine, str string) bool {
	state := 0

	for _, char := range str {
		var ok bool
		state, ok = sm.StateTransitions[state][char]
		if !ok {
			return false
		}
	}

	return sm.AcceptingStates[state]
}

func Match(sm compiler.StateMachine, str string) (ans string, matchFound bool) {
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
			fmt.Println(Test(sm, os.Args[3]))
		} else {
			fmt.Println("Invalid regex:", os.Args[2])
		}
	case "match":
		if sm, err := Compile(os.Args[2]); err == nil {
			if ans, ok := Match(sm, os.Args[3]); ok {
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

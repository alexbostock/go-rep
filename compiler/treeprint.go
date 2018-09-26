package compiler

import "fmt"

func (tree ast) print() {
	traverse(tree, 0)
}

func traverse(tree ast, indentLevel int) {
	for i := 0; i < indentLevel; i++ {
		fmt.Print("\t")
	}

	fmt.Println(tree)

	for _, child := range tree.children {
		traverse(*child, indentLevel+1)
	}
}

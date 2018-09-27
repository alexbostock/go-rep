package compiler

import (
	"errors"
	"unicode"
)

type ast struct {
	label    astlabel // the type of node
	literal  rune     // a literal value to match, applicable for literal nodes
	children []*ast   // child nodes, see below
}

// Empty nodes (only for nesting) may have many children
// Literals should have 0 children
// [xyz] constructs have as many children as runes within the []
// All other nodes should have 1 child

type astlabel int

const (
	empty astlabel = iota
	literal
	plus
	star
	optional
	anyof
	dot
)

type token struct {
	label   tokentype
	num_val int
	lit_val rune
	err     error
}

type tokentype int

const (
	tok_literal tokentype = iota
	tok_integer
	tok_openBracket
	tok_closeBracket
	tok_openSquare
	tok_closeSquare
	tok_openBrace
	tok_closeBrace
	tok_star
	tok_plus
	tok_opt
	tok_dollar
	tok_caret
	tok_dot
)

func Parse(reStr string) (syntaxTree ast, err error) {
	tokens := make(chan token)

	go lex(reStr, tokens)

	return buildTree(tokens)
}

func buildTree(tokens chan token) (syntaxTree ast, err error) {
	syntaxTree = *new(ast)
	syntaxTree.children = make([]*ast, 0, 1)

	var parsingAnyOf bool
	var parsingNumOccurences bool

	for {
		lookahead, more := <-tokens
		if !more {
			break
		}

		if parsingNumOccurences && lookahead.label != tok_integer {
			return syntaxTree, errors.New("Invalid regex: {} should contain an integer value only")
		}
		switch lookahead.label {
		case tok_dollar:
			fallthrough
		case tok_caret:
			lookahead.label = tok_literal
			lookahead.lit_val = '\n'

			fallthrough
		case tok_literal:
			newNode := new(ast)
			newNode.label = literal
			newNode.literal = lookahead.lit_val
			if parsingAnyOf {
				parent := syntaxTree.children[len(syntaxTree.children)-1]
				parent.children = append(parent.children, newNode)
			} else {
				syntaxTree.children = append(syntaxTree.children, newNode)
			}
		case tok_dot:
			newNode := new(ast)
			newNode.label = dot
			newNode.children = make([]*ast, 0)

			syntaxTree.children = append(syntaxTree.children, newNode)
		case tok_integer:
			if parsingNumOccurences {
				n := syntaxTree.children[len(syntaxTree.children)-1]
				for i := 0; i < lookahead.num_val; i++ {
					newNode := new(ast)
					newNode.label = n.label
					newNode.literal = n.literal
					newNode.children = n.children
					syntaxTree.children = append(syntaxTree.children, newNode)
				}
			} else {
				err = errors.New("Invalid regex: { should be matched with }p")
				return
			}
		case tok_openBracket:
			newNode, err := buildTree(tokens)
			if err != nil {
				return syntaxTree, err
			}
			syntaxTree.children = append(syntaxTree.children, &newNode)
		case tok_closeBracket:
			if parsingAnyOf || parsingNumOccurences {
				return syntaxTree, errors.New("Invalid regex")
			} else {
				return syntaxTree, nil
			}
		case tok_openSquare:
			parsingAnyOf = true
		case tok_closeSquare:
			parsingAnyOf = false
		case tok_openBrace:
			parsingNumOccurences = true
		case tok_closeBrace:
			parsingNumOccurences = false
		case tok_star:
			fallthrough
		case tok_plus:
			fallthrough
		case tok_opt:
			if len(syntaxTree.children) == 0 {
				return syntaxTree, errors.New("Invalid regex")
			}
			childNode := syntaxTree.children[len(syntaxTree.children)-1]
			newNode := new(ast)
			newNode.children = make([]*ast, 0, 1)
			newNode.children = append(newNode.children, childNode)
			switch lookahead.label {
			case tok_star:
				newNode.label = star
			case tok_plus:
				newNode.label = plus
			case tok_opt:
				newNode.label = optional
			}
			syntaxTree.children[len(syntaxTree.children)-1] = newNode
		}
	}

	return
}

var escapedTokens = map[rune]rune{
	'.':  '.',
	'+':  '+',
	'*':  '*',
	'{':  '{',
	'}':  '}',
	'[':  '[',
	']':  ']',
	'^':  '^',
	'$':  '$',
	'\\': '\\',
}

var tokentypes = map[rune]tokentype{
	'(': tok_openBracket,
	')': tok_closeBracket,
	'[': tok_openSquare,
	']': tok_closeSquare,
	'{': tok_openBrace,
	'}': tok_closeBrace,
	'*': tok_star,
	'+': tok_plus,
	'?': tok_opt,
	'$': tok_dollar,
	'^': tok_caret,
}

func lex(reStr string, tokens chan token) {
	var lookahead rune
	var lexingInt bool
	var lexedInt int

	// Add an EOF character
	reStr = reStr + " "

	for _, c := range reStr {
		if lookahead == 0 {
			lookahead = c
			continue
		}

		if lexingInt {
			if unicode.IsDigit(lookahead) {
				lexedInt *= 10
				lexedInt += int(lookahead)
				lookahead = c
				continue
			}
			tokens <- token{tok_integer, lexedInt, ' ', nil}
			lexingInt = false
		}

		if lookahead == '\\' {
			char, ok := escapedTokens[c]
			if !ok {
				tokens <- token{tok_literal, 0, char, errors.New("Invalid token: \\" + string(c))}
				close(tokens)
				return
			}

			tokens <- token{tok_literal, 0, char, nil}
		} else if lookahead == '.' {
			tokens <- token{tok_dot, 0, ' ', nil}
		} else if unicode.IsDigit(lookahead) {
			lexingInt = true
			lexedInt = int(lookahead)
		} else {
			label, ok := tokentypes[lookahead]
			if ok {
				tokens <- token{label, 0, ' ', nil}
			} else {
				tokens <- token{tok_literal, 0, lookahead, nil}
			}
		}

		lookahead = c
	}

	close(tokens)
}

package main

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
)

type token struct {
	label   tokentype
	num_val int
	lit_val rune
	err     error
}

type tokentype int

const (
	literal tokentype = iota
	integer
	openBracket
	closeBracket
	openSquare
	closeSquare
	openBrace
	closeBrace
	star
	plus
	opt
	dollar
	caret
)

func parse(reStr string) (syntaxTree ast, err error) {
	tokens := make(chan rune)

	go lex(reStr, tokens)

	lookahead := <-tokens

	return buildTree(lookahead, tokens)
}

func buildTree(lookahead rune, tokens chan rune) (syntaxTree ast, err error) {
	ast = new(syntaxTree)
	ast.children = make([]*ast, 1)

	var parsingAnyOf bool
	var parsingNumOccurences bool

	for token := range tokens {
		if parsingNumOccurences && lookahead.label != integer {
			return ast, errors.Nil("Invalid regex: {} should contain an integer value only")
		}
		switch lookahead.label {
		case literal:
			newNode := new(node)
			newNode.label = literal
			newNode.literal = token.lit_val
			if parsingAnyOf {
				parent := ast.children[len(ast.children)-1]
				parent.children = append(parent.children, newNode)
			} else {
				ast.children = append(ast.children, newNode)
			}
		case integer:
			if parsingNumOccurences {
				n = ast.children[len(ast.children)-1]
				for i = 0; i < token.num_val; i++ {
					newNode := new(n)
					newNode.label = n.label
					newNode.literal = n.literal
					newNode.children = n.children
					ast.children = append(ast.children, newNode)
				}
			} else {
				err = errors.New("Invalid regex: { should be matched with }")
				return
			}
		case openBracket:
			lookahead = <-tokens
			newNode, err := buildTree(lookahead, tokens)
			if err {
				return
			}
			ast.children = append(ast.children, newNode)
			lookahead = <-tokens
		case closeBracket:
			if parsingAnyOf || parsingNumOccurences {
				return ast, errors.New("Invalid regex")
			} else {
				return ast, nil
			}
		case openSquare:
			parsingAnyOf = true
		case closeSquare:
			parsingAnyOf = false
		case openBrace:
			parsingNumOccurences = true
		case closeBrace:
			parsingNumOccurences = false
		case star:
			fallthrough
		case plus:
			fallthrough
		case opt:
			if len(ast.children) == 0 {
				return ast, errors.New("Invalid regex")
			}
			childNode := ast.children[len(ast.children)-1]
			newNode := new(node)
			newNode.children = make([]*ast, 1)
			newNode.children = append(newNode.children, childNode)
			switch lookahead.label {
			case star:
				newNode.label = star
			case plus:
				newNode.label = plus
			case opt:
				newNode.label = optional
			}
			ast.children[len(ast.children)-1] = newNode
		case dollar:
			fallthrough
		case caret:
			lookahead.label = literal
			lookahead.literal = '\n'
			continue
		}

		lookahead = <-tokens
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
	'(': openBracket,
	')': closeBracket,
	'[': openSquare,
	']': closeSquare,
	'{': openBrace,
	'}': closeBrace,
	'*': star,
	'+': plus,
	'?': opt,
	'$': dollar,
	'^': caret,
}

func lex(reStr string, tokens chan token) {
	var lookahead rune
	var lexingInt bool
	var lexedInt int

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
			tokens <- token{integer, lexedInt, ' ', nil}
			lexingInt = false
		}

		if lookahead == '\\' {
			char, ok := escapedTokens[c]
			if !ok {
				tokens <- token{literal, 0, char, errors.New("Invalid token: \\" + string(c))}
				close(tokens)
				return
			}

			tokens <- token{literal, 0, char, nil}
		} else if unicode.IsDigit(lookahead) {
			lexingInt = true
			lexedInt = int(lookahead)
		} else {
			label, ok := tokentypes[lookahead]
			if ok {
				tokens <- token{label, 0, ' ', nil}
			} else {
				tokens <- token{literal, 0, lookahead, nil}
			}
		}

		lookahead = c
	}

	close(tokens)
}

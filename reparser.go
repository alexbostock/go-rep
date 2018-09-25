package main

import (
	"errors"
	"unicode"
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

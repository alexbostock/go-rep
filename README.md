# go-rep
A regular expression processor, just for fun.

## Usage

`StateMachine struct { ... }`: represents a finite state machine to accept a regular language.

`Compile(reStr string) (sm StateMachine, err error)`: parse a regex string `reStr` and create an equivalent state machine `sm`.

`(sm StateMachine) Test(str string) bool`: returns true iff `sm` accepts `str`.

`(sm StateMachine) Match(str string) (ans string, matchFound bool)`: if `str` contains a string accepted by `sm`, return the first such matching string as `ans`. `matchFound` is true iff any match is found.

## Syntax
- `.`: any single character
- `x+`: 1 or more occurrences of `x`
- `x*`: 0 or more occurrences of `x`
- `x?`: 0 or 1 occurrences of `x`
- `x{n}`: `n` occurrences of `x`
- `^`: beginning of string or beginning of line
- `$`: end of string or end of line
- `[xyz]`: any one of `x`, `y`, `z`
- `[a-z]`: any character whose character codes lies between the codes of `a` and `z`, inclusive
- `\.`, `\+`, `\*`, `\{`, `\}`, `\[`, `\]`, `\^`, `\$`, `\`: special characters
- `(`, `)`: changing precedence

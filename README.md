# go-rep
A regular expression processor, just for fun.

## Usage

`Compile`: parse a regex string and generate a state machine.

`Test`: return true iff a given string is accepted by a given state machine.

`Match`: return a string which is the first match of a given regex (state machine) in a given string.

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

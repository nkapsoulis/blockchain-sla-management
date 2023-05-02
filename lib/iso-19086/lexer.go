package iso19086

type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Parameters
	PARAM

	// Variables
	VAR

	// Constants
	CONST

	// Operators
	ADD
	SUB
	MUL
	DIV
	LESS
	MORE
)

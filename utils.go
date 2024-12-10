package main

import "regexp"

// define regular expressions

var (
	TokenRegex      = regexp.MustCompile(`([a-zA-Z\d_$#@%^;:~.""*+\-/'\"]+)|([=])|([{}])|([()])|([|])|([&])|(['\"]|([<>])|([!])|([?]))`)
	BracketsRegex   = regexp.MustCompile(`[()]`)
	OperRegex       = regexp.MustCompile(`[+-/*]`)
	AndORRegex      = regexp.MustCompile(`[&|]`)
	LogicalRegex    = regexp.MustCompile(`[!<>=]`)
	DecimalRegex    = regexp.MustCompile(`^[\d\.]*$`)
	FloatRegex      = regexp.MustCompile(`^\d*[\.]?\d+$`)
	IntRegex        = regexp.MustCompile(`^\d+$`)
	DecimalExpRegex = regexp.MustCompile(`^[\d\.][+-\/*%=<>!]*[\d]+`)
	VariableRegex   = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z_\d]*$`)
	CombinedRegex   = regexp.MustCompile(`[a-zA-Z\d_\.+-\/*()]`)
)

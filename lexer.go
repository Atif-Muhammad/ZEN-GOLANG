package main

import (
	"regexp"
	"strings"
)

func Lexer(line_num int, line string, datatypes []string, grammarSymbols []string) (string, error) {

	// clean and tokenize each line
	content := strings.Trim(line, "\t\n ")

	tokens := []string{}
	
	tkns := TokenRegex.FindAllString(content, -1)
	
	/* ignore comments */
	if tkns[0] == "|" && tkns[len(tkns)-1] == "|" {
	} else if len(tkns) > 1 || (len(tkns) == 1 && tkns[0] == "else") {
		if includes(datatypes, tkns[0]) {
			// prepare the tokens according to the pre defined rules for each dataType
			// e.g.
			// datatype var = value || datatype var = 2+1 || datatype var = var+2 || datatype var = var+var || datatype var = get('string')
			// it means the length of the value is not known, so we will combine all the tokens after the '=' symbol
			last_token := strings.Join(tkns[3:], " ")
			tokens = append(tokens, tkns[0], tkns[1], tkns[2], last_token)
		} else if includes(grammarSymbols, tkns[0]) {
			last_token := strings.Join(tkns[2:], " ")
			tokens = append(tokens, tkns[0], tkns[1], last_token)
		} else if tkns[0] == "if" {
			if tkns[len(tkns)-1] == "{" {
				exp := strings.Join(tkns[2:len(tkns)-2], " ")
				tokens = append(tokens, tkns[0], tkns[1], exp, tkns[len(tkns)-2], tkns[len(tkns)-1])
			} else if tkns[len(tkns)-1] == ")" {
				exp := strings.Join(tkns[2:len(tkns)-1], " ")
				tokens = append(tokens, tkns[0], tkns[1], exp, tkns[len(tkns)-1])
			}
		} else if tkns[0] == "else" || tkns[0] == "}" {
			tokens = append(tokens, tkns...)
		} else if DecimalRegex.MatchString(tkns[0]) || VariableRegex.MatchString(tkns[0]) {
			if len(tkns) <= 5 {
				logical_operator := strings.Join(tkns[1:len(tkns)-2], " ")
				tokens = append(tokens, tkns[0], logical_operator, tkns[len(tkns)-2], tkns[len(tkns)-1])
			} else {
				state := strings.Join(tkns[0:len(tkns)-1], " ")
				tokens = append(tokens, state, tkns[len(tkns)-1])
			}
		}
	} else if len(tkns) == 1 {
		match, _ := regexp.MatchString(tkns[0], "[{}]")
		if match {
			tokens = append(tokens, tkns...)
		} else {
			// throw error
			err := ErrorHandler(line_num, line, "Syntax Error", "Invalid Syntax")
			if err != nil {
				return "", err;
			}
		}
	}

	// call parser for tokens
	var final_result string
	if len(tokens) > 0{
		result, err := Parser(line_num, line, tokens, datatypes, grammarSymbols);
		if(err != nil){
			return "", err
		}else{
			final_result = result
		}
	}
	return final_result, nil
}

package main

import (
	// "fmt"
	"strings"
)

var stack = []string{}

func checkMatchingBracket(bracket string) bool {

	brackets := map[string]string{"{": "}"}
	if bracket == "{" {
		stack = append(stack, bracket)
	} else if bracket == "}" {
		
		if len(stack) == 0 {
			return false
		}
		// pop from the stack and check it matches the closing bracket?
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if brackets[top] != bracket{
			// fmt.Println("After pop: ", stack)
			return false
		}
	}
	return true
}

func Parser(line_num int, line string, tokens []string, datatypes []string, grammarSymbols []string) (string, error){

	var final_result = ""

	if len(tokens) > 1 || (len(tokens) == 1 && tokens[0] == "else"){
		if includes(datatypes, tokens[0]){
			if tokens[2] == "="{
				if strings.TrimSpace(tokens[3]) != ""{
					// pass tokens to semantic analyzer
					result, err := SemanticAnalyzer(line_num ,line ,tokens ,datatypes ,grammarSymbols);
					if(err != nil){
						return "", err
					}else{
						final_result = result;
					}
				}else{
					// throw error
					err := ErrorHandler(line_num, line, "Value Error", "Missing Value")
					if err != nil {
						return "", err;
					}
				}
			}else{
				err := ErrorHandler(line_num, line, "Syntax Error", "Invalid Syntax")
				if err != nil {
					return "", err;
				}
			}		
		}else if includes(grammarSymbols, tokens[0]){
			if tokens[0] == "Set"{
				// Set>"string...."
				if tokens[1] == ">"{
					result, err := SemanticAnalyzer(line_num ,line ,tokens ,datatypes ,grammarSymbols)
					if(err != nil){
						return "", err
					}else{
						final_result = result;
					}
				}else{
					err := ErrorHandler(line_num, line, "Syntax Error", "Invalid Syntax")
					if err != nil {
						return "", err;
					}
				}
			}else{
				err := ErrorHandler(line_num, line, "Syntax Error", "Invalid Syntax for Get")
				if err != nil {
					return "", err;
				}
			}
		}else if tokens[0] == "if"{
			if tokens[len(tokens)-1] == "{"{
				checkMatchingBracket(tokens[len(tokens)-1])
				if tokens[1] == "("{
					if tokens[len(tokens)-2] == ")"{
						result, err := SemanticAnalyzer(line_num ,line ,tokens ,datatypes ,grammarSymbols)
						if(err != nil){
							return "", err
						}else{
							final_result = result;
						}
					}else{
						err := ErrorHandler(line_num, line, "Syntax Error", "Invalid Syntax")
						if err != nil {
							return "", err;
						}
					}
				}else{
					err := ErrorHandler(line_num, line, "Syntax Error", "Invalid Syntax")
					if err != nil {
						return "", err;
					}
				}
				
			}else if tokens[len(tokens)-1] == ")"{
				if tokens[1] == "("{
					result, err := SemanticAnalyzer(line_num ,line ,tokens ,datatypes ,grammarSymbols)
					if(err != nil){
						return "", err
					}else{
						final_result = result;
					}
				}else{
					err := ErrorHandler(line_num, line, "Syntax Error", "Invalid Syntax")
					if err != nil {
						return "", err;
					}
				}
			}
		}else if tokens[0] == "else" || tokens[0] == "}"{
			if (tokens[len(tokens)-1] != "" && tokens[len(tokens)-1] == "{"){
				checkMatchingBracket(tokens[len(tokens)-1]);
			}
			if tokens[0] == "else"{
				result, err := SemanticAnalyzer(line_num ,line ,tokens ,datatypes ,grammarSymbols)
				if(err != nil){
					return "", err
				}else{
					final_result = result;
				}
			}else if tokens[0] == "}"{
				if(checkMatchingBracket(tokens[0])){
					if tokens[1] == "else"{
						result, err := SemanticAnalyzer(line_num ,line ,tokens ,datatypes ,grammarSymbols)
						if(err != nil){
							return "", err
						}else{
							final_result = result;
						}
					}else{
						err := ErrorHandler(line_num, line, "Syntax Error", "Invalid Syntax")
						if err != nil {
							return "", err;
						}
					}
				}else{
					err := ErrorHandler(line_num, line, "Syntax Error", "Bracket Error")
					if err != nil {
						return "", err;
					}
				}
			}
		}else if tokens[len(tokens)-1] == "?"{
			if len(tokens) == 2{
				result, err := SemanticAnalyzer(line_num ,line ,tokens ,datatypes ,grammarSymbols)
				if(err != nil){
					return "", err
				}else{
					final_result = result;
				}
			}else if len(tokens) > 2{
				if DecimalRegex.MatchString(tokens[len(tokens)-2]) || VariableRegex.MatchString(tokens[len(tokens)-2]){
					// pass to semantics
					result, err := SemanticAnalyzer(line_num ,line ,tokens ,datatypes ,grammarSymbols)
					if(err != nil){
						return "", err
					}else{
						final_result = result;
					}
				}
			}
		}else{
			err := ErrorHandler(line_num, line, "Syntax Error", "Invalid Syntax")
			if err != nil {
				return "", err;
			}
		}
	}else if len(tokens) == 1 {	
		if tokens[0] == "{" || tokens[0] == "}"{
			if(tokens[0] == "{"){
				checkMatchingBracket(tokens[0]);
			}else{
				if(checkMatchingBracket(tokens[0])){
					// pass to semantic
					result, err := SemanticAnalyzer(line_num ,line ,tokens ,datatypes ,grammarSymbols)
					if(err != nil){
						return "", err
					}else{
						final_result = result;
					}
				}else{
					err := ErrorHandler(line_num, line, "Syntax Error", "Bracket Error")
					if err != nil {
						return "", err;
					}
				}
			}
		}else{
			// throw error invalid syntax
			err := ErrorHandler(line_num, line, "Syntax Error", "Invalid Syntax")
			if err != nil {
				return "", err;
			}
		}
	}

	return final_result, nil;
}

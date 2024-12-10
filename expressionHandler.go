
package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/expr-lang/expr"
)

// define important global variables for return
var IntReturn *int
var FloatReturn *float64
var BoolReturn *bool

func getReturn() interface{} {
	switch {
	case IntReturn != nil:
		return *IntReturn
	case FloatReturn != nil:
		return *FloatReturn
	case BoolReturn != nil:
		return *BoolReturn
	default:
		return nil
	}
}

func ExpressionHandler(line_num int, line string, exp string) (interface{}, error) {



	FloatReturn = nil
	IntReturn = nil
	BoolReturn = nil
	
	if DecimalRegex.MatchString(exp) {
		if IntRegex.MatchString(exp) {
			// convert string to int
			num, err := strconv.Atoi(exp)
			if err != nil {
				return "", err
								
			} else {
				IntReturn = &num
			}
		} else if FloatRegex.MatchString(exp) {
			// convert string to float
			num, err := strconv.ParseFloat(exp, 64)
			if err != nil {
				return "", err
				
			} else {
				FloatReturn = &num
			}
		} else {
		}
	} else if exp == "true" || exp == "false" {
		retval, _ := strconv.ParseBool(exp)
		BoolReturn = &retval
	} else if VariableRegex.MatchString(exp) {
		// retrieve value from symbol table and return it
		val, err := RetrieveST(line_num, line, exp)
		if(err != nil){
			return "", err
		}else{
			if value, ok := val.(int); ok {
				IntReturn = &value
			} else if value, ok := val.(float64); ok {
				FloatReturn = &value
			} else if value, ok := val.(bool); ok {
				BoolReturn = &value
			}
		}
	} else if CombinedRegex.MatchString(exp) {
		var temp_exp string = ""
		var expression string = ""
		if LogicalRegex.MatchString(exp) {
			// fmt.Println(exp)
			for _, char := range exp {
				if LogicalRegex.MatchString(string(char)) {
					if VariableRegex.MatchString(temp_exp) {
						var value string
						if temp_exp == "true" || temp_exp == "false" {
							value = temp_exp
						} else {
							// retrieve value from symbol table
							val, err := RetrieveST(line_num, line, temp_exp)
							if(err != nil){
								return "", err
							}else{
								value = fmt.Sprintf("%v", val)
							}
						}
						temp_exp = ""
						expression += value
						expression += string(char)
					} else {
						expression += temp_exp
						temp_exp = ""
						expression += string(char)
					}
				} else if BracketsRegex.MatchString(string(char)) {

				} else if char == ' ' {

				} else {
					temp_exp += string(char)
				}
			}
			if temp_exp != "" {
				if VariableRegex.MatchString(temp_exp) {
					var value string
					if temp_exp == "true" || temp_exp == "false" {
						value = temp_exp
					} else {
						// retrieve value from symbol table
						val, err := RetrieveST(line_num, line, temp_exp)
						if(err != nil){
							return "", err
						}else{
							value = fmt.Sprintf("%v", val)
						}
					}
					temp_exp = ""
					expression += value
				} else {
					expression += temp_exp
					temp_exp = ""
				}
			}
			if expression != "" {
				// call evaluator to evaluate the value or expression
				// 	and update the BoolReturn value
				result, err := expr.Eval(expression, nil)
				if err != nil {
					return "", err
					
				} else {
					if val, ok := result.(bool); ok {
						BoolReturn = &val
					} else {
						return "", fmt.Errorf("Unknownn Expression Error")
						
					}
				}
			}
		} else if AndORRegex.MatchString(exp) {
			for _, char := range exp {
				if AndORRegex.MatchString(string(char)) {
					if temp_exp == "true" || temp_exp == "false" {
						expression += temp_exp
						expression += " "
						expression += string(char)
						temp_exp = ""
					} else {
						expression += string(char)
					}
				} else {
					temp_exp += string(char)
				}
			}

			if temp_exp != "" {
				if temp_exp == "true" || temp_exp == "false" {
					expression += " "
					expression += temp_exp
					temp_exp = ""
				} else {
					// throw error-unexpected
					return "", fmt.Errorf("Unknownn Expression Error")
					
				}
			}

			if expression != "" {
				// call evaluator to evaluate the value or expression
				// 	and update the BoolReturn value
				AlphaRegex := regexp.MustCompile(`[A-Za-z]`)
				var updated_exp string
				for _, char := range expression {
					if AlphaRegex.MatchString(string(char)) {
						updated_exp += string(char)
					} else if char == '&' {
						updated_exp += " && "
					} else if char == '|' {
						updated_exp += " || "
					}
				}
				result, err := expr.Eval(updated_exp, nil)
				if err != nil {
					return "", err
					
				} else {
					if val, ok := result.(bool); ok {
						BoolReturn = &val
					} else {
						return "", fmt.Errorf("Unknownn Expression Error")
						
					}
				}
			}
		} else {
			for _, char := range exp {
				if OperRegex.MatchString(string(char)) {
					if DecimalRegex.MatchString(temp_exp) {
						expression += temp_exp
						expression += string(char)
						temp_exp = ""
					} else if VariableRegex.MatchString(temp_exp) {
						// retrieve value from symbol table
						val, err := RetrieveST(line_num, line, temp_exp)
						if(err != nil){
							return "", err
						}else{
							value := fmt.Sprintf("%v", val)
							temp_exp = ""
							expression += value
							expression += string(char)
						}
					}
				} else if BracketsRegex.MatchString(string(char)) {

				} else if char == ' ' {

				} else if char == '.' {
					temp_exp += string(char)
				} else {
					temp_exp += string(char)
				}
			}

			if temp_exp != "" {
				if DecimalRegex.MatchString(temp_exp) {
					expression += temp_exp
					temp_exp = ""
				} else if VariableRegex.MatchString(temp_exp) {
					// retrieve value from symbol table
					val, err := RetrieveST(line_num, line, temp_exp)
					if(err != nil){
						return "", err
					}else{
						value := fmt.Sprintf("%v", val)
						temp_exp = ""
						expression += value
					}
				}
			}

			if expression != "" {
				// evaluate the expression and return the result
				res, err := expr.Eval(expression, nil)
				if err != nil {
					return "", err
					
				} else {
					if result, ok := res.(int); ok {
						IntReturn = &result
					} else if result, ok := res.(float64); ok {
						FloatReturn = &result
					} else {
						return "", fmt.Errorf("Unknownn Expression Error")
						
					}
				}
			}
		}
	} else if LogicalRegex.MatchString(exp) {
		
		result, err := expr.Eval(exp, nil)
		if err != nil {
			return "", err
			
		} else {
			if res, ok := result.(bool); ok {
				BoolReturn = &res
			}
		}
	} 
	return getReturn(), nil
}

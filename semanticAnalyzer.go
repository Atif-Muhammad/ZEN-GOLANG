package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var Executing bool = true
var exec_block []bool
var if_count int = 0

func SemanticAnalyzer(line_num int, line string, tkns []string, datatypes []string, grammarSymbols []string) (string, error) {

	tokens := []string{}

	var result = ""

	// variable for removing quote signs from the string assigned to Literal
	var final_string string
	// analyzed tokens for multiple types i.e., int, float, string, bool
	var analyzed_tokens []interface{}

	// clean the tokens
	for i := 0; i < len(tkns); i++ {
		tokens = append(tokens, strings.Trim(tkns[i], " "))
	}

	if tokens[0] == "{" && len(tokens) == 1 {
	} else if tokens[0] == "}" && len(tokens) == 1 {
		if if_count > 0 {
			if len(exec_block) > 0 {

				if if_count == 1 {
					Executing = true
				} else {
					if exec_block[if_count-2] {
						Executing = exec_block[len(exec_block)-1]
						exec_block = exec_block[:len(exec_block)-1]
					}
				}
				if_count = if_count - 1
			}
		} else {
			Executing = true
		}
	} else if includes(datatypes, tokens[0]) && Executing {
		// validate the variable name
		if VariableRegex.MatchString(tokens[1]) {
			// The value(after '=' sign) can either be the value, variable, expression or a statement
			// check if a statement
			if tokens[0] == "Num" {
				r, _ := regexp.Compile("[<]")
				if r.MatchString(tokens[len(tokens)-1]) {
					// pass the last token to Get method to get user input
					/* check whether user input is of int type
					if yes: append it to the analyzed_tokens and pass them to symbol table
					if no: throw error of invalid assignment to Num
					*/
					user_in, err := Get(line_num, line, tokens[len(tokens)-1])
					if err != nil {
						return "", err
					} else {
						if IntRegex.MatchString(user_in) {
							rtrn, err := strconv.Atoi(user_in)
							if err != nil {
								return "", err
							} else {
								analyzed_tokens = append(analyzed_tokens, tokens[0], tokens[1], rtrn)
							}
						} else {
							err := ErrorHandler(line_num, line, "Value Error", "Invalid Assignment to Num")
							if err != nil {
								return "", err
							}
						}
					}

				} else if IntRegex.MatchString(tokens[len(tokens)-1]) {
					// parse the last token into int
					num, err := strconv.Atoi(tokens[len(tokens)-1])
					if err != nil {
						// fmt.Println(err)
					} else {
						analyzed_tokens = append(analyzed_tokens, tokens[0], tokens[1], num)
					}
				} else if VariableRegex.MatchString(tokens[len(tokens)-1]) {
					/* retrieve variable value from symbol table
					and check whether it is an int or not(error)*/

					value, err := RetrieveST(line_num, line, tokens[len(tokens)-1])
					if err != nil {
						return "", err
					} else {
						if reflect.TypeOf(value).Kind() == reflect.Int {
							analyzed_tokens = append(analyzed_tokens, tokens[0], tokens[1], value)

						} else {
							err := ErrorHandler(line_num, line, "Value Error", "Invalid Assignment to Num")
							if err != nil {
								return "", err
							}
						}
					}
				} else if CombinedRegex.MatchString(tokens[len(tokens)-1]) {
					// check if an expression and pass it to expression handler
					/*	check if the expression's return is an integer
						if yes: append it to analyzed_tokens and pass them to symbol table
						if no: throw error of invalid assignment to Num
					*/
					result, err := ExpressionHandler(line_num, line, tokens[len(tokens)-1])
					if err != nil {
						return "", err
					} else {
						if val, ok := result.(int); ok {
							analyzed_tokens = append(analyzed_tokens, tokens[0], tokens[1], val)
						} else {
							err := ErrorHandler(line_num, line, "Value Error", "Invalid Assignment to Num")
							if err != nil {
								return "", err
							}
						}
					}
				} else {
					// throw invalid assignment to Num error-value
					err := ErrorHandler(line_num, line, "Value Error", "Invalid Assignment to Num")
					if err != nil {
						return "", err
					}
				}
			} else if tokens[0] == "Decimal" {
				r, _ := regexp.Compile("[<]")
				if r.MatchString(tokens[len(tokens)-1]) {
					/* check whether user input is of float type
					if yes: append it to the analyzed_tokens and pass them to symbol table
					if no: throw error of invalid assignment to Decimal
					*/
					user_in, err := Get(line_num, line, tokens[len(tokens)-1])
					if err != nil {
						return "", err
					} else {
						if FloatRegex.MatchString(user_in) {
							rtrn, err := strconv.ParseFloat(user_in, 64)
							if err != nil {
								// throw err
							} else {
								analyzed_tokens = append(analyzed_tokens, tokens[0], tokens[1], rtrn)
							}
						} else {
							err := ErrorHandler(line_num, line, "Value Error", "Invalid Assignment to Decimal")
							if err != nil {
								return "", err
							}
						}
					}

				} else if FloatRegex.MatchString(tokens[len(tokens)-1]) {
					// parse the last token to int
					num, err := strconv.ParseFloat(tokens[len(tokens)-1], 64)
					if err != nil {
						// fmt.Println(err)
					} else {
						analyzed_tokens = append(analyzed_tokens, tokens[0], tokens[1], num)
					}
				} else if VariableRegex.MatchString(tokens[len(tokens)-1]) {
					/* retrieve variable value from symbol table
					and check whether it is an float or not(error)*/
					value, err := RetrieveST(line_num, line, tokens[len(tokens)-1])
					if err != nil {
						return "", err
					} else {
						if reflect.TypeOf(value).Kind() == reflect.Float64 {
							analyzed_tokens = append(analyzed_tokens, tokens[0], tokens[1], value)
						} else {
							err := ErrorHandler(line_num, line, "Value Error", "Invalid Assignment to Decimal")
							if err != nil {
								return "", err
							}
						}
					}
				} else if CombinedRegex.MatchString(tokens[len(tokens)-1]) {
					// check if an expression and pass it to expression handler
					/*	check if the expression's return is a float
						if yes: append it to analyzed_tokens and pass them to symbol table
						if no: throw error of invalid assignment to Decimal
					*/
					result, err := ExpressionHandler(line_num, line, tokens[len(tokens)-1])
					if err != nil {
						return "", err
					} else {
						if val, ok := result.(float64); ok {
							analyzed_tokens = append(analyzed_tokens, tokens[0], tokens[1], val)
						} else {
							err := ErrorHandler(line_num, line, "Value Error", "Invalid Assignment to Decimal")
							if err != nil {
								return "", err
							}
						}
					}
				} else {
					// invalid assignment to Decimal
					err := ErrorHandler(line_num, line, "Value Error", "Invalid Assignment to Decimal")
					if err != nil {
						return "", err
					}
				}
			} else if tokens[0] == "Literal" {
				r, _ := regexp.Compile("[<]")
				if r.MatchString(tokens[len(tokens)-1]) {
					// pass the last token to Get method to get user input
					/* check whether user input is of string type
					if yes: append it to the analyzed_tokens and pass them to symbol table
					if no: throw error of invalid assignment to Literal
					*/
					user_in, err := Get(line_num, line, tokens[len(tokens)-1])
					if err != nil {
						return "", err
					} else {
						analyzed_tokens = append(analyzed_tokens, tokens[0], tokens[1], user_in)
					}
				} else if VariableRegex.MatchString(tokens[len(tokens)-1]) {
					// retrieve value from symbol table and append into the analyzed_tokens slice
					value, err := RetrieveST(line_num, line, tokens[len(tokens)-1])
					if err != nil {
						return "", err
					} else {
						if reflect.TypeOf(value).Kind() == reflect.String {
							analyzed_tokens = append(analyzed_tokens, tokens[0], tokens[1], value)
						} else {
							err := ErrorHandler(line_num, line, "Value Error", "Invalid Assignment to Literal")
							if err != nil {
								return "", err
							}
						}
					}
				} else if (strings.HasPrefix(tokens[len(tokens)-1], `"`) && strings.HasSuffix(tokens[len(tokens)-1], `"`)) || (strings.HasPrefix(tokens[len(tokens)-1], "'") && strings.HasSuffix(tokens[len(tokens)-1], "'")) {
					// remove the quote sign from the string
					if strings.HasPrefix(tokens[len(tokens)-1], `"`) && strings.HasSuffix(tokens[len(tokens)-1], `"`) {
						_string := strings.ReplaceAll(tokens[len(tokens)-1], `"`, " ")
						final_string = strings.Trim(_string, " ")
					} else if strings.HasPrefix(tokens[len(tokens)-1], "'") && strings.HasSuffix(tokens[len(tokens)-1], "'") {
						_string := strings.ReplaceAll(tokens[len(tokens)-1], `'`, " ")
						final_string = strings.Trim(_string, " ")
					} else {
						// throw error for miss match quotes
						err := ErrorHandler(line_num, line, "Syntax Error", "Mismatching Quotes")
						if err != nil {
							return "", err
						}
					}
					analyzed_tokens = append(analyzed_tokens, tokens[0], tokens[1], final_string)
				}
			} else if tokens[0] == "Bool" {
				r1, _ := regexp.Compile(`(^[a-zA-Z_][a-zA-Z_\d]*$)|\d+\s*(<=|>=|==|!=|<|>)\s*(^[a-zA-Z_][a-zA-Z_\d]*$)|\d+`)
				r2, _ := regexp.Compile(`[<]`)
				if r2.MatchString(tokens[len(tokens)-1]) {
					/*
						pass the last token to get the user input and check whether it is either
							true or false and then append them into the expression handler
					*/
					user_in, err := Get(line_num, line, tokens[len(tokens)-1])
					if err != nil {
						return "", err
					} else {
						if user_in == "true" || user_in == "false" {
							rtrn, err := strconv.ParseBool(user_in)
							if err != nil {
								// throw err
							} else {
								analyzed_tokens = append(analyzed_tokens, tokens[0], tokens[1], rtrn)
							}
						} else {
							err := ErrorHandler(line_num, line, "Value Error", "Invalid Assignment to Bool")
							if err != nil {
								return "", err
							}
						}
					}
				} else if tokens[len(tokens)-1] == "true" || tokens[len(tokens)-1] == "false" {
					if tokens[len(tokens)-1] == "true" {
						analyzed_tokens = append(analyzed_tokens, tokens[0], tokens[1], true)
					} else {
						analyzed_tokens = append(analyzed_tokens, tokens[0], tokens[1], false)
					}
				} else if VariableRegex.MatchString(tokens[len(tokens)-1]) {
					/*
						retrieve variable from symbol table
							check whether the value is bool????
					*/
					value, err := RetrieveST(line_num, line, tokens[len(tokens)-1])
					if err != nil {
						return "", err
					} else {
						if reflect.TypeOf(value).Kind() == reflect.Bool {
							analyzed_tokens = append(analyzed_tokens, tokens[0], tokens[1], value)
						} else {
							err := ErrorHandler(line_num, line, "Value Error", "Invalid Assignment to Bool")
							if err != nil {
								return "", err
							}
						}
					}
				} else if r1.MatchString(tokens[len(tokens)-1]) {
					// i.e., var == var
					// call the expression handler to resolve the variables and return bool values
					// append the tokens with the result from expression handler
					result, err := ExpressionHandler(line_num, line, tokens[len(tokens)-1])
					if err != nil {
						return "", err
					} else {
						if val, ok := result.(bool); ok {
							analyzed_tokens = append(analyzed_tokens, tokens[0], tokens[1], val)
						} else {
							err := ErrorHandler(line_num, line, "Value Error", "Invalid Assignment to Bool")
							if err != nil {
								return "", err
							}
						}
					}
				} else {
					// throw error
					err := ErrorHandler(line_num, line, "Value Error", "Invalid Assignment to Bool")
					if err != nil {
						return "", err
					}
				}
			}
		} else {
			// throw error
			err := ErrorHandler(line_num, line, "Syntax Error", "Invalid Variable Name")
			if err != nil {
				return "", err
			}
		}
		// pass the anlyzed tokens to Symbol Table
		err := CreateST(line_num, line, analyzed_tokens)
		if err != nil {
			return "", err
		}
	} else if includes(grammarSymbols, tokens[0]) && Executing {
		if tokens[0] == "Set" {
			// declare essential variables and flags for better evaluation of strings/variables
			var variable bool = false
			var variable_name string = ""
			var temp_message string = ""
			var quote_sign rune
			var in_quotes bool = false
			var quote_sign_count int = 0

			if strings.HasPrefix(tokens[len(tokens)-1], `"`) || strings.HasPrefix(tokens[len(tokens)-1], `'`) {
				r3, _ := regexp.Compile(`[+-/*]`)
				// if string with expression or variable name
				for _, char := range tokens[len(tokens)-1] {
					if (char == '"' || char == '\'') && !in_quotes {
						in_quotes = true
						quote_sign = char
						temp_message += string(char)
						quote_sign_count = quote_sign_count + 1
					} else if char != quote_sign && in_quotes {
						temp_message += string(char)
					} else if char == quote_sign && in_quotes {
						in_quotes = false
						temp_message += string(char)
						quote_sign_count = quote_sign_count + 1
					} else if char == '>' && !in_quotes && !variable {
						variable = true
					} else if char != ' ' && variable {
						variable_name += string(char)
					} else if char != ' ' && !variable {
						// throw error unexpected syntax
						err := ErrorHandler(line_num, line, "Syntax Error", "Unexpected Syntax")
						if err != nil {
							return "", err
						}
					} else if r3.MatchString(string(char)) {
						variable_name += string(char)
					}
				}
				if quote_sign_count < 2 {
					err := ErrorHandler(line_num, line, "Syntax Error", "Missing Quotation")
					if err != nil {
						return "", err
					}
				} else if quote_sign_count > 2 {
					err := ErrorHandler(line_num, line, "Syntax Error", "Quotation Error")
					if err != nil {
						return "", err
					}
				}
			} else {
				// if expression or variable only
				err := Set(line_num, line, "", "", tokens[len(tokens)-1])
				if err != nil {
					return "", err
				}
			}

			if temp_message != "" || variable_name != "" {
				if variable_name == "" && variable {
					// throw error missing value
					err := ErrorHandler(line_num, line, "Value Error", "Missing Value")
					if err != nil {
						return "", err
					}
				} else {
					if temp_message != "" {
						// remove quote signs
						msg := strings.Trim(strings.ReplaceAll(temp_message, string(quote_sign), ""), "")
						// sent the msg to Set along with the variable name
						err := Set(line_num, line, msg, variable_name, "")
						if err != nil {
							return "", err
						}
					} else {
						err := Set(line_num, line, temp_message, variable_name, "")
						if err != nil {
							return "", err
						}
					}
				}
			}
		} else {
			// do nothing for Get<"......"
			// or throw error
			err := ErrorHandler(line_num, line, "Syntax Error", "Invalid Syntax for Get")
			if err != nil {
				return "", err
			}
		}
	} else if tokens[0] == "if" {
		if_count += 1
		if strings.Count(tokens[0], "&") > 1 || strings.Count(tokens[0], "|") > 1 {
			err := ErrorHandler(line_num, line, "Logic Error", "Unexpected Logic")
			if err != nil {
				return "", err
			}
		} else {
			if tokens[len(tokens)-1] == "{" || tokens[len(tokens)-1] == ")" {
				if Executing {
					if AndORRegex.MatchString(tokens[2]) {
						var temp_exp string = ""
						var expression string = ""
						for _, char := range tokens[2] {
							if AndORRegex.MatchString(string(char)) {
								if LogicalRegex.MatchString(temp_exp) {
									// call expression handler for bool value
									result, err := ExpressionHandler(line_num, line, temp_exp)
									if err != nil {
										return "", err
									} else {
										temp_exp = ""
										expression += fmt.Sprint(result)
										expression += string(char)
									}

								} else {
									expression += string(char)
								}

							} else if char == ' ' {

							} else {
								temp_exp += string(char)
							}
						}

						if temp_exp != "" {
							if LogicalRegex.MatchString(temp_exp) {
								// call expression handler
								result, err := ExpressionHandler(line_num, line, temp_exp)
								if err != nil {
									return "", err
								} else {
									temp_exp = ""
									expression += fmt.Sprint(result)
								}
							} else {
								expression += temp_exp
							}
						}

						if expression != "" {
							// call expression handler
							result, err := ExpressionHandler(line_num, line, expression)
							if err != nil {
								return "", err
							} else {
								if res, ok := result.(bool); ok {
									if Executing {
										Executing = res
									}
									exec_block = append(exec_block, res)
								}
							}
						}
					} else {
						if LogicalRegex.MatchString(tokens[2]) {
							// call the expression handler for boolean value
							// and update the Executing value
							result, err := ExpressionHandler(line_num, line, tokens[2])
							if err != nil {
								return "", err
							} else {
								if res, ok := result.(bool); ok {
									if Executing {
										Executing = res
									}
									exec_block = append(exec_block, res)
								}
							}
						} else if VariableRegex.MatchString(tokens[2]) {
							// call symbol table
							// update the executing value accordingly
							value, err := RetrieveST(line_num, line, tokens[2])
							if err != nil {
								return "", err
							} else {
								if result, ok := value.(bool); ok {
									if Executing {
										Executing = result
									}
									exec_block = append(exec_block, result)
								}
							}
						}
					}
				}
			}
		}
	} else if tokens[0] == "else" || tokens[0] == "}" {

		if if_count > 0 {
			if len(exec_block) > 0 {
				if if_count == 1 {
					Executing = !exec_block[0]
				} else if if_count > 2 {
					if exec_block[if_count-2] {
						Executing = exec_block[if_count-1]
					}
				} else {
					if exec_block[if_count-2] {
						Executing = !exec_block[if_count-1]
					} else {
						Executing = false
					}
				}
			}
		} else {
			err := ErrorHandler(line_num, line, "Syntax Error", "Invalid Syntax for else")
			if err != nil {
				return "", err
			}
		}
	} else if tokens[len(tokens)-1] == "?" && Executing {
		if len(tokens) == 2 {
			if strings.Count(tokens[0], "&") > 1 || strings.Count(tokens[0], "|") > 1 {
				// throw error
				err := ErrorHandler(line_num, line, "Logic Error", "Unexpected Logic")
				if err != nil {
					return "", err
				}
			} else {
				var temp_exp string = ""
				var expression string = ""
				for _, char := range tokens[0] {
					if AndORRegex.MatchString(string(char)) {
						if LogicalRegex.MatchString(temp_exp) {
							// call expression handler for bool value
							result, err := ExpressionHandler(line_num, line, temp_exp)
							if err != nil {
								return "", err
							} else {
								temp_exp = ""
								expression += fmt.Sprint(result)
								expression += string(char)
							}

						} else {
							expression += string(char)
						}

					} else if char == ' ' {

					} else {
						temp_exp += string(char)
					}
				}

				if temp_exp != "" {
					if LogicalRegex.MatchString(temp_exp) {
						// call expression handler
						result, err := ExpressionHandler(line_num, line, temp_exp)
						if err != nil {
							return "", err
						} else {
							temp_exp = ""
							expression += fmt.Sprint(result)
						}
					} else {
						expression += temp_exp
					}
				}

				if expression != "" {
					// call expression handler
					result, err := ExpressionHandler(line_num, line, expression)
					if err != nil {
						return "", err
					} else {
						res := fmt.Sprintf("%v", result)
						temp_exp = ""
						fmt.Println(res)
						// result = res
						// return res, nil
					}

				}
			}
		} else {
			if DecimalRegex.MatchString(tokens[0]) {
				if DecimalRegex.MatchString(tokens[len(tokens)-2]) {
					analyzed_tokens = append(analyzed_tokens, tokens[0], tokens[1], tokens[len(tokens)-2])
				} else if VariableRegex.MatchString(tokens[len(tokens)-2]) {
					// retrieve data from symbol table
					if tokens[len(tokens)-2] == "true" {
						analyzed_tokens = append(analyzed_tokens, tokens[0], tokens[1], true)
					} else if tokens[len(tokens)-2] == "false" {
						analyzed_tokens = append(analyzed_tokens, tokens[0], tokens[1], false)
					} else {
						value, err := RetrieveST(line_num, line, tokens[len(tokens)-2])
						if err != nil {
							return "", err
						} else {
							analyzed_tokens = append(analyzed_tokens, tokens[0], tokens[1], value)
						}
					}
				} else {
					// throw missing value error
					err := ErrorHandler(line_num, line, "Value Error", "Missing Value")
					if err != nil {
						return "", err
					}
				}
			} else if VariableRegex.MatchString(tokens[0]) {
				// retrieve tokens[0] and tokens[last]
				var value1 interface{}
				var value2 interface{}

				if tokens[0] == "true" {
					value1 = true
				} else if tokens[0] == "false" {
					value1 = false
				} else {
					value1, _ = RetrieveST(line_num, line, tokens[0])
					_, err := RetrieveST(line_num, line, tokens[0])
					if err != nil {
						return "", err
					}
				}
				if DecimalRegex.MatchString(tokens[len(tokens)-2]) {
					analyzed_tokens = append(analyzed_tokens, value1, tokens[1], tokens[2])
				} else if VariableRegex.MatchString(tokens[len(tokens)-2]) {
					if tokens[len(tokens)-2] == "true" {
						value2 = true
					} else if tokens[len(tokens)-2] == "false" {
						value2 = false
					} else {
						value2, _ = RetrieveST(line_num, line, tokens[len(tokens)-2])
						_, err := RetrieveST(line_num, line, tokens[len(tokens)-2])
						if err != nil {
							return "", err
						}

					}
					analyzed_tokens = append(analyzed_tokens, value1, tokens[1], value2)
				} else {
					// throw error missing value
					err := ErrorHandler(line_num, line, "Value Error", "Missing Value")
					if err != nil {
						return "", err
					}
				}
			}

			if analyzed_tokens != nil {
				// convert each element to a string
				var strItem []string
				for _, item := range analyzed_tokens {
					strItem = append(strItem, fmt.Sprint(item))
				}
				exp := strings.Join(strItem, " ")
				// call expression handler for exp
				result, err := ExpressionHandler(line_num, line, exp)
				if err != nil {
					return "", err
				} else {
					res := fmt.Sprintf("%v", result)
					// return res, nil
					fmt.Println(res)
				}

			}
		}
	}
	return result, nil
}

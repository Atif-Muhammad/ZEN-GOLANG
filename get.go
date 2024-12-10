package main

import (
	"fmt"
	"strings"
)

func Get(line_num int, line string, content string) (string, error){

	// declare return
	var user_in string
	// tokenize the content
	tkns := TokenRegex.FindAllString(content, -1)

	if tkns[0] == "Get"{
		if tkns[1] == "<"{
			msg := strings.Join(tkns[2:], " ")
			if (strings.HasPrefix(msg, `'`) && strings.HasSuffix(msg, `'`)) || (strings.HasPrefix(msg, `"`) && strings.HasSuffix(msg, `"`)){
				in_quotes := false
				var quote_sign rune
				temp_message := ""

				for _, char := range msg{
					if (char == '"' || char == '\'') && !in_quotes{
						quote_sign = char
						in_quotes = true
					}else if char != quote_sign && in_quotes{
						temp_message+=string(char)
					}else if char == quote_sign && in_quotes{
						in_quotes = false
					}
				}

				if strings.HasPrefix(msg, string(quote_sign)) && strings.HasSuffix(msg, string(quote_sign)){
					// accept user input
					fmt.Print(temp_message, " ")
					_, err := fmt.Scan(&user_in)
					if err != nil{
						return "", err
					}

				}else{
					// throw quote error
					err := ErrorHandler(line_num, line, "Syntax Error", "Quote Error")
					return "", err
				}				
			}else{
				// throw error for quote
				err := ErrorHandler(line_num, line, "Value Error", "Unexpected Error")
				return "", err
			}
		}else{
			// throw missing < error
			err := ErrorHandler(line_num, line, "Syntax Error", "Missing <")
			return "", err
		}
	}else{
		// throw error
		err := ErrorHandler(line_num, line, "Name Error", "Unexpected Syntax")
		return "", err
	}

	return user_in, nil
}
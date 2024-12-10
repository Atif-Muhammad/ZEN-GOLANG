package main

import (
	"fmt"
)

func Set(line_num int, line string, message string, variable string, exp string) error {

	// if message and variable
	if variable != "" && message != "" {
		value, err := ExpressionHandler(line_num, line, variable)
		if(err != nil){
			return err
		}else{
			fmt.Println(message , value)
		}
	} else if message != "" && variable == "" && exp == "" {
		// if only message
		fmt.Println(message)
	} else if variable != "" {
		// if only variable
		value, err := ExpressionHandler(line_num, line, variable)
		if(err != nil){
			return err
		}else{
			fmt.Println(value)
		}
	}

	if exp != "" {
		if LogicalRegex.MatchString(exp){
			value, err := ExpressionHandler(line_num, line, exp)
			if(err != nil){
				return err
			}else{
				fmt.Println(value)
			}
		}else{
			value, err := RetrieveST(line_num, line, exp)
			if(err != nil){
				return err
			}else{
				fmt.Println(value)
			}
		}
	}
	return nil
}

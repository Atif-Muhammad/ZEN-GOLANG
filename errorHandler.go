package main

import (
	"fmt"
)

func ErrorHandler(line_num int, line string, errorType string, desc string) error {
	
	error_msg := fmt.Sprintf("%s: %s in line '%d' at [%s]", errorType, desc, line_num, line)
	
	return fmt.Errorf("%s", error_msg)
}

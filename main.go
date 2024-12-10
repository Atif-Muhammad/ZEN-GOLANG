package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

/* define include function
returns true if item is present in the slice
returns false if item is not present in the slice */

func includes(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func main() {

	// Define dataTypes
	datatypes := []string{"Num", "Decimal", "Literal", "Bool"}

	// Define grammarSymbols
	grammarSymbols := []string{"Get", "Set"}

	/////////////////////////////////////////////////////////////////////////////
	/* to accept user input from command line - external file */
	if len(os.Args) < 2 {
		fmt.Println("Error: no file path provided.")
		os.Exit(1)
	}

	// Get the file path from command-line
	filepath := os.Args[1]
	// open the file
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	// declare line tracker
	var line_num int = 0
	response := ""
	for scanner.Scan() {
		line := scanner.Text()
		// pass each line to Lexer
		line_num += 1
		if strings.TrimSpace(line) != "" {
			result, err := Lexer(line_num, line, datatypes, grammarSymbols)
			if err != nil {
				// panic(err)
				log.Println(err)
				break
			} else {
				if result != "" {
					response += result + "\n"
				}

			}
		}
	}

	log.Println(response)

}

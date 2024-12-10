package main


type symbol struct {
	dataType interface{}
	name     interface{}
	value    interface{}
}

// global map to store variables
var symbol_table_dict = make(map[interface{}]symbol)

func CreateST(line_num int, line string, tokens []interface{}) error {

	dict_name := tokens[1]

	_, ok := symbol_table_dict[dict_name]
	if ok {
		err := ErrorHandler(line_num, line, "Value Error", "Cannot Initialize Existing Variable")
		if err != nil {
			return err
		}
	}

	symbol_table_dict[dict_name] = symbol{tokens[0], dict_name, tokens[len(tokens)-1]}
	return nil
}

func RetrieveST(line_num int, line string, variable string) (interface{}, error) {

	_, ok := symbol_table_dict[variable]

	if ok {
		// fmt.Println(symbol_table_dict[variable].value)
		val := symbol_table_dict[variable].value
		return val, nil
	} else {
		err := ErrorHandler(line_num, line, "Value Error", "Accessing Before Initialization")
		return "", err
	}
}

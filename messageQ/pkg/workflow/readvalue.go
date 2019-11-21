package workflow

import (
	"encoding/json"
	"errors"
	"strings"
)

func getStringFromJSONMap(value string, JSONMap map[string]string) (string, error) {

	//Get the first value from the dictionary
	nameList := strings.Split(value, ".")
	first, present := JSONMap[nameList[0]]
	if !present {
		return "", errors.New("No value '" + first + "' found")
	}
	if len(nameList) == 1 {
		return first, nil
	}

	// and make it a map[string]*json.RawMessage
	var objmap map[string]*json.RawMessage
	err := json.Unmarshal([]byte(first), &objmap)
	if err != nil {
		return "", err
	}

	// get the JSON value and make it a new map until the last variable name is reached
	var result *json.RawMessage
	nameList = nameList[1:]
	for index, name := range nameList {
		result, present = objmap[name]
		if !present {
			return "", errors.New("No value '" + value + "' found")
		}
		if index < len(nameList)-1 {
			objmap = nil
			if err := json.Unmarshal(*result, &objmap); err != nil {
				return "", err
			}
		}
	}

	//Convert the last value to a string if found
	finalValue := ""
	if result != nil {
		finalValue = string(*result)
		if len(finalValue) > 2 {
			finalValue = finalValue[1 : len(finalValue)-1] //Supress the double quote
		}
	}

	return finalValue, nil
}

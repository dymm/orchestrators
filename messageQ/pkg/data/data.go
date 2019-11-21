package data

import (
	"encoding/json"
	"errors"
)

//TestValue is the test data type
type TestValue struct {
	Name  string
	Value int
}

//DeserializeTestValue is deserializing the test value
func DeserializeTestValue(values map[string]string) (TestValue, error) {

	serialized, found := values["data"]

	var data TestValue
	if !found {
		return data, errors.New("No data found")
	}

	err := json.Unmarshal([]byte(serialized), &data)
	return data, err
}

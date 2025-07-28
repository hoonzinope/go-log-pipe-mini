package parse

import (
	"encoding/json"
	"fmt"
	"test_gluent_mini/shared"
)

func ParseJSON(input string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(input), &result)
	if err != nil {
		shared.Error_count.Add(1) // Increment error count
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}
	return result, nil
}

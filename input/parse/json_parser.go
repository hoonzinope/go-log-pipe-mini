package parse

import (
	"encoding/json"
	"fmt"
)

func ParseJSON(input string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(input), &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}
	return result, nil
}

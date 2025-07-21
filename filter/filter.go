package filter

import (
	"strings"
)

func _grep(line string, filterIgnoreCase bool, filterPattern string) bool {
	keywords := strings.Split(filterPattern, "|") // Split the pattern by pipe character
	flag := false                                 // Flag to indicate if the line matches any keyword
	for _, keyword := range keywords {
		if filterIgnoreCase {
			if strings.Contains(strings.ToLower(line), strings.ToLower(keyword)) {
				flag = true // Set flag to true if the line contains the keyword (case insensitive)
				break       // Break out of the loop if a match is found
			}
		} else {
			if strings.Contains(line, keyword) {
				flag = true // Set flag to true if the line contains the keyword (case sensitive)
				break       // Break out of the loop if a match is found
			}
		}
	}
	return flag // Return true if the line matches any keyword, false otherwise
}

func _json_grep(jsonLine map[string]interface{}, field string,
	filterIgnoreCase bool, filterPattern string) bool {
	if jsonLine == nil || field == "" || filterPattern == "" {
		return false
	}
	flag := false                                 // Flag to indicate if the line matches the filter
	keywords := strings.Split(filterPattern, "|") // Split the pattern by pipe character
	for _, keyword := range keywords {
		if value, exists := jsonLine[field]; exists {
			if strValue, ok := value.(string); ok {
				if filterIgnoreCase {
					if strings.Contains(strings.ToLower(strValue), strings.ToLower(keyword)) {
						flag = true // Set flag to true if the field contains the pattern (case insensitive)
						break
					}
				} else {
					if strings.Contains(strValue, keyword) {
						flag = true // Set flag to true if the field contains the pattern (case sensitive)
						break
					}
				}
			}
		}
	}
	return flag // Return true if the field matches the pattern, false otherwise
}

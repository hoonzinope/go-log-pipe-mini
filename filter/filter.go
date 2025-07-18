package filter

import (
	"context"
	"strings"
	"test_gluent_mini/confmanager"
	"test_gluent_mini/data"
)

var _mode string = "OR" // Default mode, can be changed based on config
var _funcList = make([]func(data.InputData) bool, 0)

func Configure(conf confmanager.Config) {
	_mode = strings.ToUpper(conf.Filter.Mode) // Set the mode from the configuration
	for _, filter := range conf.Filter.Filters {
		f := filter // Create a local copy of the filter to avoid closure issues
		switch f.Type {
		case "grep":
			// Register the grep function with its options
			_funcList = append(_funcList, func(line data.InputData) bool {
				return _grep(line.Raw, f.Options.IgnoreCase, f.Options.Pattern)
			})
		case "json_grep":
			_funcList = append(_funcList, func(line data.InputData) bool {
				return _json_grep(line.Json, f.Options.Field, f.Options.IgnoreCase, f.Options.Pattern)
			})
		default:
			// Handle other filter types if needed
			continue
		}
	}
}

func FilterLine(ctx context.Context, logLineChan chan data.InputData, filterLineChan chan string) {
	for {
		select {
		case <-ctx.Done():
			return // Exit if the context is cancelled
		case lineInputData := <-logLineChan:
			if filterFunc(lineInputData) {
				filterLineChan <- lineInputData.Raw
			}
		}
	}
}

func filterFunc(lineData data.InputData) bool {
	flag := false // Initialize flag to false

	if _mode == "OR" {
		for _, f := range _funcList {
			if f(lineData) {
				flag = true // Set flag to true if any filter function returns true
				break
			}
		}
		return flag // Return true if any filter function matched, false otherwise
	} else if _mode == "AND" {
		flag = true // Start with true for AND mode
		for _, f := range _funcList {
			if !f(lineData) {
				flag = false // Set flag to false if any filter function returns false
				break
			}
		}
		return flag // Return true only if all filter functions matched
	} else {
		panic("Unknown filter mode: " + _mode) // Panic if the mode is unknown
	}
}

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

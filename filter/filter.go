package filter

import (
	"context"
	"strings"
	"test_gluent_mini/confmanager"
)

var _mode string = "OR" // Default mode, can be changed based on config
var _funcList = make([]func(string) bool, 0)

func Configure(conf confmanager.Config) {
	_mode = strings.ToUpper(conf.Filter.Mode) // Set the mode from the configuration
	for _, filter := range conf.Filter.Filters {
		f := filter // Create a local copy of the filter to avoid closure issues
		switch f.Type {
		case "grep":
			// Register the grep function with its options
			_funcList = append(_funcList, func(line string) bool {
				return _grep(line, f.Options.IgnoreCase, f.Options.Pattern)
			})
		default:
			// Handle other filter types if needed
			continue
		}
	}
}

func FilterLine(ctx context.Context, logLineChan chan string, filterLineChan chan string) {
	for {
		select {
		case <-ctx.Done():
			return // Exit if the context is cancelled
		case line := <-logLineChan:
			if filterFunc(line) {
				filterLineChan <- line
			}
		}
	}
}

func filterFunc(line string) bool {
	flag := false // Initialize flag to false

	if _mode == "OR" {
		for _, f := range _funcList {
			if f(line) {
				flag = true // Set flag to true if any filter function returns true
				break
			}
		}
		return flag // Return true if any filter function matched, false otherwise
	} else if _mode == "AND" {
		flag = true // Start with true for AND mode
		for _, f := range _funcList {
			if !f(line) {
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
	var flag bool = false                         // Flag to indicate if the line matches any keyword
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

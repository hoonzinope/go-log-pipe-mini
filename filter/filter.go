package filter

import (
	"context"
	"strings"
	"test_gluent_mini/confmanager"
)
var filterPattern string // Global variable to hold the filter pattern
var filterIgnoreCase bool // Global variable to hold the ignore case option
var filterFunc func(string) bool = _grep // Default filter function

func Configure(conf confmanager.Config) {
	filterType := conf.Filter.Type
	switch filterType {
	case "":
		filterType = "grep" // Default to grep if no filter type is specified
		filterPattern = conf.Filter.Options.Patterns // Set the filter pattern from configuration
		filterIgnoreCase = conf.Filter.Options.Ignore_Case // Set the ignore case option from configuration
		filterFunc = _grep // Set the filter function to grep
	case "grep":
		filterType = "grep" // Set filter type to grep
		filterPattern = conf.Filter.Options.Patterns // Set the filter pattern from configuration
		filterIgnoreCase = conf.Filter.Options.Ignore_Case // Set the ignore case option from configuration
		filterFunc = _grep // Set the filter function to grep
	default:
		filterType = "grep" // Default to grep if unsupported filter type is specified
		filterPattern = conf.Filter.Options.Patterns // Set the filter pattern from configuration
		filterIgnoreCase = conf.Filter.Options.Ignore_Case // Set the ignore case option from configuration
		filterFunc = _grep // Set the filter function to grep
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

func _grep(line string) bool {
	keywords := strings.Split(filterPattern, "|") // Split the pattern by pipe character
	var flag bool = false // Flag to indicate if the line matches any keyword
	for _, keyword := range keywords {
		if filterIgnoreCase {
			if strings.Contains(strings.ToLower(line), strings.ToLower(keyword)) {
				flag = true // Set flag to true if the line contains the keyword (case insensitive)
				break // Break out of the loop if a match is found
			}
		} else {
			if strings.Contains(line, keyword) {
				flag = true // Set flag to true if the line contains the keyword (case sensitive)
				break // Break out of the loop if a match is found
			}
		}
	}
	return flag // Return true if the line matches any keyword, false otherwise
}
package filter

import (
	"context"
	"strings"
)

func FilterLine(ctx context.Context, logLineChan chan string, filterLineChan chan string, keyword string) {
	for {
		select {
		case <-ctx.Done():
			return // Exit if the context is cancelled
		case line := <-logLineChan:
			if _grep(line, keyword) {
				filterLineChan <- line
			}
		}
	}
}

func _grep(line string, keyword string) bool {
	if keyword == "" {
		return true // If no keyword is specified, return true to pass all lines
	}
	return strings.Contains(line, keyword) || strings.Contains(line, strings.ToLower(keyword))
	// You can also use regex or other complex filtering logic here if needed
}
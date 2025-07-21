package filter

import (
	"context"
	"strings"
	"test_gluent_mini/confmanager"
	"test_gluent_mini/data"
)

var cancel_ctx context.Context
var config confmanager.Config
var logLineChan map[string]chan data.InputData
var filterLineChan map[string]chan string

func Configure(ctx context.Context,
	conf confmanager.Config,
	inputChannel map[string]chan data.InputData,
	filterChannel map[string]chan string,
) {
	cancel_ctx = ctx
	config = conf
	logLineChan = inputChannel
	filterLineChan = filterChannel
}

func FilterLines() {
	var filters map[string]confmanager.FilterGroup = config.Filters
	for name, filterGroup := range filters {
		go _runFilterGroup(cancel_ctx, name, filterGroup)
	}
}

func _runFilterGroup(ctx context.Context,
	name string, filterGroup confmanager.FilterGroup) {
	_mode := strings.ToUpper(filterGroup.Mode)
	var _funcList = make([]func(data.InputData) bool, 0)
	for _, rules := range filterGroup.Rules {
		rule := rules // Create a local copy of the rule to avoid closure issues
		switch rule.Type {
		case "grep":
			// Register the grep function with its options
			_funcList = append(_funcList, func(line data.InputData) bool {
				return _grep(line.Raw, rule.Options.IgnoreCase, rule.Options.Pattern)
			})
		case "json_grep":
			_funcList = append(_funcList, func(line data.InputData) bool {
				return _json_grep(line.Json, rule.Options.Field, rule.Options.IgnoreCase, rule.Options.Pattern)
			})
		default:
			// Handle other filter types if needed
			continue
		}
	}
	_filterLine(ctx, name, _mode, _funcList)
}

func _filterLine(ctx context.Context,
	name string, mode string, funcList []func(data.InputData) bool) {
	for {
		select {
		case <-ctx.Done():
			return // Exit if the context is cancelled
		case lineInputData := <-logLineChan[name]:
			if filterFunc(lineInputData, mode, funcList) {
				filterLineChan[name] <- lineInputData.Raw
			}
		}
	}
}

func filterFunc(lineData data.InputData, mode string, funcList []func(data.InputData) bool) bool {
	flag := false // Initialize flag to false

	if mode == "OR" {
		for _, f := range funcList {
			if f(lineData) {
				flag = true // Set flag to true if any filter function returns true
				break
			}
		}
		return flag // Return true if any filter function matched, false otherwise
	} else if mode == "AND" {
		flag = true // Start with true for AND mode
		for _, f := range funcList {
			if !f(lineData) {
				flag = false // Set flag to false if any filter function returns false
				break
			}
		}
		return flag // Return true only if all filter functions matched
	} else {
		return false // Unsupported mode, return false
	}
}

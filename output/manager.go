package output

import (
	"context"
	"fmt"
	"test_gluent_mini/confmanager"
)

var cancel_ctx context.Context
var config confmanager.Config
var filterLineChan map[string]chan string

func Configure(
	ctx context.Context,
	conf confmanager.Config,
	filterLineChannel map[string]chan string) {
	config = conf
	cancel_ctx = ctx
	filterLineChan = filterLineChannel
}

func Out() {
	for _, outputConfig := range config.Outputs {
		switch outputConfig.Type {
		case "stdout":
			// Print to standard output
			var outputFunc = _println
			for _, target := range outputConfig.Targets {
				go _out(cancel_ctx, outputFunc, filterLineChan[target])
			}
		default:
			fmt.Printf("Unsupported output type: %s\n", outputConfig.Type)
		}
	}
}

func _out(ctx context.Context, outputFunc func(string), lineChan chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case logLine := <-lineChan:
			if logLine != "" {
				outputFunc(logLine)
			}
		}
	}
}

func _println(logLine string) {
	fmt.Printf("%s\n", logLine)
}

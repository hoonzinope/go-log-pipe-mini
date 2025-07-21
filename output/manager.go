package output

import (
	"context"
	"fmt"
	"test_gluent_mini/confmanager"
	"test_gluent_mini/shared"
)

var config confmanager.Config

func Configure(conf confmanager.Config) {
	config = conf
}

func Out() {
	for _, outputConfig := range config.Outputs {
		switch outputConfig.Type {
		case "stdout":
			// Print to standard output
			var outputFunc = _println
			for _, target := range outputConfig.Targets {
				go _out(shared.Ctx, outputFunc, shared.FilterChannel[target])
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

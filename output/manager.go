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

func _out(ctx context.Context, outputFunc func(string), lineChan chan shared.InputData) {
	for {
		select {
		case <-ctx.Done():
			return
		case logLine := <-lineChan:
			if logLine.Json != nil {
				outputFunc(fmt.Sprintf("%s %s: %v", logLine.Tag, logLine.FileName, logLine.Json))
			} else if logLine.Raw != "" {
				outputFunc(fmt.Sprintf("%s %s: %s", logLine.Tag, logLine.FileName, logLine.Raw))
			}
			shared.OffsetChannel <- logLine
		}
	}
}

func _println(logLine string) {
	fmt.Printf("%s\n", logLine)
}

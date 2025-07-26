package output

import (
	"context"
	"fmt"
	"slices"
	"test_gluent_mini/confmanager"
	"test_gluent_mini/shared"
)

const (
	ROLLING_DEFAULT        = "daily"
	MAX_SIZE_DEFAULT       = "100MB"
	MAX_FILES_DEFAULT      = 7
	BATCH_SIZE_DEFAULT     = 10
	FLUSH_INTERVAL_DEFAULT = "1s"
)

var config confmanager.Config

type Outputer interface {
	Out(ctx context.Context)
}

func Configure(conf confmanager.Config) {
	config = conf
}

func Out() {

	outputChannel := make(map[string]map[string]chan shared.InputData, len(config.Outputs))
	for _, outputConfig := range config.Outputs {
		if _, exists := outputChannel[outputConfig.Type]; !exists {
			outputChannel[outputConfig.Type] = make(map[string]chan shared.InputData)
		}
		for _, target := range outputConfig.Targets {
			ch := make(chan shared.InputData, 1000)
			outputChannel[outputConfig.Type][target] = ch
		}
	}

	for _, outputConfig := range config.Outputs {
		switch outputConfig.Type {
		case "stdout":
			consoleOutput := ConsoleOutput{
				Type:           outputConfig.Type,
				Targets:        outputConfig.Targets,
				BATCH_SIZE:     outputConfig.Options.BATCH_SIZE,
				FLUSH_INTERVAL: outputConfig.Options.FLUSH_INTERVAL,
			}
			consoleOutput.Out(shared.Ctx, outputChannel[outputConfig.Type])
		case "file":
			fileOutput := FileOutput{
				Targets:        outputConfig.Targets,
				Path:           outputConfig.Options.Path,
				Filename:       outputConfig.Options.Filename,
				Rolling:        outputConfig.Options.Rolling,
				MaxSize:        outputConfig.Options.MaxSize,
				MaxFiles:       outputConfig.Options.MaxFiles,
				BATCH_SIZE:     outputConfig.Options.BATCH_SIZE,
				FLUSH_INTERVAL: outputConfig.Options.FLUSH_INTERVAL,
			}
			fileOutput.Out(shared.Ctx, outputChannel[outputConfig.Type])
		case "http":
			httpOutput := HttpOutput{
				Type:           outputConfig.Type,
				Targets:        outputConfig.Targets,
				Url:            outputConfig.Options.Url,
				Method:         outputConfig.Options.Method,
				Headers:        outputConfig.Options.Headers,
				Timeout:        outputConfig.Options.Timeout,
				BATCH_SIZE:     outputConfig.Options.BATCH_SIZE,
				FLUSH_INTERVAL: outputConfig.Options.FLUSH_INTERVAL,
			}
			httpOutput.Out(shared.Ctx, outputChannel[outputConfig.Type])
		default:
			fmt.Printf("Unsupported output type: %s\n", outputConfig.Type)
		}
	}

	go _broadcastFilteredData(outputChannel)
}

func _broadcastFilteredData(outputChannel map[string]map[string]chan shared.InputData) {
	targets := make([]string, 0, len(outputChannel))
	for _, outputConfig := range config.Outputs {
		for _, target := range outputConfig.Targets {
			if !slices.Contains(targets, target) {
				targets = append(targets, target)
			}
		}
	}
	for {
		select {
		case <-shared.Ctx.Done():
			return // Exit if the context is cancelled
		default:
			for _, target := range targets {
				select {
				case <-shared.Ctx.Done():
					return // Exit if the context is cancelled
				case logLine := <-shared.InputChannel[target]:
					for outputType, channels := range outputChannel {
						if ch, exists := channels[target]; exists {
							ch <- logLine
						} else {
							fmt.Printf("No channel found for target %s in output type %s\n", target, outputType)
						}
					}
					shared.OffsetChannel <- logLine // send to offset channel after broadcasting
				default:
					// No data to broadcast, continue
				}
			}
		}
	}
}

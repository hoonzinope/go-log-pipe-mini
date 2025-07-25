package output

import (
	"context"
	"fmt"
	"test_gluent_mini/confmanager"
	"test_gluent_mini/shared"
)

const (
	ROLLING_DEFAULT      = "daily"
	MAX_SIZE_DEFAULT     = "100MB"
	MAX_FILES_DEFAULT    = 7
	BATCH_SIZE_DEFAULT   = 10
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
	for _, outputConfig := range config.Outputs {
		switch outputConfig.Type {
		case "stdout":
			consoleOutput := ConsoleOutput{
				Type:    outputConfig.Type,
				Targets: outputConfig.Targets,
				BATCH_SIZE: outputConfig.Options.BATCH_SIZE,
				FLUSH_INTERVAL: outputConfig.Options.FLUSH_INTERVAL,
			}
			consoleOutput.Out(shared.Ctx)
		case "file":
			fileOutput := FileOutput{
				Targets:  outputConfig.Targets,
				Path:     outputConfig.Options.Path,
				Filename: outputConfig.Options.Filename,
				Rolling:  outputConfig.Options.Rolling,
				MaxSize:  outputConfig.Options.MaxSize,
				MaxFiles: outputConfig.Options.MaxFiles,
				BATCH_SIZE: outputConfig.Options.BATCH_SIZE,
				FLUSH_INTERVAL: outputConfig.Options.FLUSH_INTERVAL,
			}
			fileOutput.Out(shared.Ctx)
		default:
			fmt.Printf("Unsupported output type: %s\n", outputConfig.Type)
		}
	}
}

package output

import (
	"context"
	"fmt"
	"test_gluent_mini/confmanager"
	"test_gluent_mini/shared"
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
			}
			fileOutput.Out(shared.Ctx)
		default:
			fmt.Printf("Unsupported output type: %s\n", outputConfig.Type)
		}
	}
}

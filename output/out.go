package output

import (
	"context"
	"fmt"
	"test_gluent_mini/confmanager"
)

var outputFunc func(string) = _println // Default output function
func Configure(conf confmanager.Config) {
	outputType := conf.Output.Type
	switch outputType {
	case "":
		outputFunc = _println // Default to Println if no output type is specified
	case "stdout":
		outputFunc = _println // Set output function to Println for standard output
	default:
		outputFunc = _println // Default to Println if unsupported output type is specified
	}

}

func Out(ctx context.Context, filterLineChannel chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case logLine := <-filterLineChannel:
			if logLine != "" {
				outputFunc(logLine) // Print to standard output
			}
		}
	}
}

func _println(logLine string) {
	fmt.Println(logLine) // Print to standard output
}

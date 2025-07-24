package output

import (
	"context"
	"fmt"
	"test_gluent_mini/shared"
)

type ConsoleOutput struct {
	Type    string
	Targets []string
}

func (c ConsoleOutput) Out(ctx context.Context) {
	// Print to standard output
	for _, target := range c.Targets {
		lineChan := shared.FilterChannel[target]
		go func(ctx context.Context, lineChan chan shared.InputData) {
			for {
				select {
				case <-ctx.Done():
					return
				case logLine := <-lineChan:
					if logLine.Json != nil {
						fmt.Printf("%s %s: %v\n", logLine.Tag, logLine.FileName, logLine.Json)
					} else if logLine.Raw != "" {
						fmt.Printf("%s %s: %s\n", logLine.Tag, logLine.FileName, logLine.Raw)
					}
					shared.OffsetChannel <- logLine
				}
			}
		}(ctx, lineChan)
	}
}

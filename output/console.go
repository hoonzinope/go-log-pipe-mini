package output

import (
	"context"
	"fmt"
	"test_gluent_mini/shared"
	"time"
)

type ConsoleOutput struct {
	Type    string
	Targets []string
	BATCH_SIZE int
	FLUSH_INTERVAL string
}

func (c ConsoleOutput) Out(ctx context.Context) {
	if c.BATCH_SIZE == 0 {
		c.BATCH_SIZE = BATCH_SIZE_DEFAULT
	}
	if c.FLUSH_INTERVAL == "" {
		c.FLUSH_INTERVAL = FLUSH_INTERVAL_DEFAULT
	}
	
	// Print to standard output
	duration, err := time.ParseDuration(c.FLUSH_INTERVAL)
	if err != nil {
		fmt.Printf("Error parsing FLUSH_INTERVAL %s: %v\n", c.FLUSH_INTERVAL, err)
		return
	}
	for _, target := range c.Targets {
		lineChan := shared.FilterChannel[target]
		go func(ctx context.Context, lineChan chan shared.InputData) {
			for {
				batch := make([]shared.InputData, 0, c.BATCH_SIZE)
				timer := time.NewTimer(duration)
				BATCHLOOP: 
				for {
					select {
					case <-ctx.Done():
						return
					case logLine := <-lineChan:
						// Write immediately if no batching or flushing is configured
						batch = append(batch, logLine)
						if len(batch) >= c.BATCH_SIZE {
							break BATCHLOOP
						}
					case <-timer.C:
						break BATCHLOOP
					}
				}
				
				for _, logLine := range batch {
					if err := c._writeToConsole(logLine); err != nil {
						fmt.Printf("Error writing to console: %v\n", err)
					}
					shared.OffsetChannel <- logLine
				}
			}
		}(ctx, lineChan)
	}
}

func (c ConsoleOutput) _writeToConsole(logLine shared.InputData) error {
	if logLine.Json != nil {
		fmt.Printf("%s %s: %v\n", logLine.Tag, logLine.FileName, logLine.Json)
	} else if logLine.Raw != "" {
		fmt.Printf("%s %s: %s\n", logLine.Tag, logLine.FileName, logLine.Raw)
	}
	return nil
}

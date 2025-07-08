package output

import (
	"context"
	"fmt"
)

func Stdout(ctx context.Context, loglineChannel chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case logLine := <-loglineChannel:
			if logLine != "" {
				println(logLine) // Print to standard output
			}
		}
	}
}

func println(logLine string) {
	fmt.Print(logLine) // Print to standard output
}
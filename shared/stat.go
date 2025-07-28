package shared

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

var Input_count atomic.Int64 = atomic.Int64{}
var Filter_count atomic.Int64 = atomic.Int64{}
var Output_count atomic.Int64 = atomic.Int64{}
var Error_count atomic.Int64 = atomic.Int64{}

var durationSum atomic.Int64 = atomic.Int64{}
var processedCount atomic.Int64 = atomic.Int64{}

func AddLatency(duration time.Duration) {
	if duration < 0 {
		return // Ignore negative durations
	}
	durationSum.Add(int64(duration))
	processedCount.Add(1)
}

func GetAverageLatency() time.Duration {
	if processedCount.Load() == 0 {
		return 0
	}
	s := durationSum.Load()
	c := processedCount.Load()
	if c == 0 {
		return 0
	}
	avg := s / c
	return time.Duration(avg)
}

func PrintStats(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(_print())
			return
		default:
			fmt.Println(_print())
			// Sleep for a while before printing again
			time.Sleep(5 * time.Second)
		}
	}
}

func _print() string {
	return fmt.Sprintf("[stat] Input: %d | Filter: %d | Output: %d | Error: %d | Avg Latency: %s",
		Input_count.Load(),
		Filter_count.Load(),
		Output_count.Load(),
		Error_count.Load(),
		GetAverageLatency())
}

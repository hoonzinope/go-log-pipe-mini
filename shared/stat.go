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

var durationSum int64 = 0
var processedCount int64 = 0

func AddLatency(duration time.Duration) {
	atomic.AddInt64(&durationSum, int64(duration))
	atomic.AddInt64(&processedCount, 1)
}

func GetAverageLatency() time.Duration {
	if processedCount == 0 {
		return 0
	}
	s := atomic.LoadInt64(&durationSum)
	c := atomic.LoadInt64(&processedCount)
	if c == 0 {
		return 0
	}
	avg := s / c
	return time.Duration(avg) * time.Nanosecond
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

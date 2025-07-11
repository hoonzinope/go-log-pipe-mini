package input

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"test_gluent_mini/offset"
	"time"
)
var lastOffset int64 = -1 // Variable to store the last offset read from the file
var lineProcessedCount int64 = 0 // Global variable to count processed lines
var lastFlushTime time.Time = time.Now() // Variable to track the last flush time

func init() {
	// Initialize the lastOffset variable by reading from the offset file
	offset, err := offset.ReadOffset() // Read the offset from the file
	if err != nil {
		fmt.Printf("Error reading offset: %v\n", err)
		lastOffset = 0 // If there's an error, start from the beginning
	} else {
		fmt.Printf("Last offset read from file: %d\n", offset)
		lastOffset = offset // Set the lastOffset to the value read from the file
	}
}

func TailFile(ctx context.Context, loglineChannel chan string, filePath string, offsetChannel chan int64) {
	offset := lastOffset // Use the last offset read from the file
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, stopping input reading.")
			return
		default:
			lastLine, lastOffset := _tailFile(filePath, offset)
			if lastLine != "" {
				lineProcessedCount++ // Increment the count of processed lines
				offset = lastOffset // Update the offset with the new value
				loglineChannel <- lastLine
				if lineProcessedCount%1000 == 0 || time.Since(lastFlushTime) > 10*time.Second {
					offsetChannel <- lastOffset // Send the updated offset to the channel
					lastFlushTime = time.Now() // Update the last flush time
				}
			}
		}
	}
}

func _tailFile(filePath string, offset int64) (string,int64) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return "", offset
	}
	defer file.Close()
	file.Seek(offset, io.SeekStart)
	line, err := bufio.NewReader(file).ReadString('\n')
	if err != nil {
		// fmt.Printf("Error reading file %s: %v %v\n", filePath, err, offset)
		time.Sleep(5 * time.Second) // Wait before retrying
		return "", offset // Return empty string if no new line is read
	}
	offset += int64(len(line)) // Update the offset for the next read
	return line, offset
}
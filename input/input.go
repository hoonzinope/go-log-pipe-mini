package input

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"time"
)
func TailFile(ctx context.Context, loglineChannel chan string, filePath string, offset int64, offsetChannel chan int64) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, stopping input reading.")
			return
		default:
			lastLine, lastOffset := _tailFile(filePath, offset)
			if lastLine != "" {
				offset = lastOffset // Update the offset with the new value
				loglineChannel <- lastLine
				offsetChannel <- lastOffset // Send the updated offset to the channel
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
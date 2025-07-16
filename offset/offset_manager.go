package offset

import (
	"context"
	"fmt"
	"os"
	"time"
)

var offsetFileTemp string = "./offset.tmp"
var offsetFilePath string = "./offset.state"
var offsetMap map[string]int64 = nil
var lastFlushTime time.Time = time.Now()

func GetOffsetMap() (map[string]int64, error) {
	if offsetMap == nil {
		offsetMap = make(map[string]int64)
		offsets, err := _read()
		if err != nil {
			fmt.Printf("Error reading offsets: %v\n", err)
		} else {
			for file, off := range offsets {
				offsetMap[file] = off
				fmt.Printf("Last offset for %s: %d\n", file, off)
			}
		}
	}
	return offsetMap, nil
}

func _read() (map[string]int64, error) {
	offsets := make(map[string]int64)
	offsetFilePath := "./offsets.state"
	file, err := os.Open(offsetFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return offsets, nil // If the file doesn't exist, return an empty map
		}
		fmt.Printf("Error opening offsets file %s: %v\n", offsetFilePath, err)
		return nil, err
	}
	defer file.Close()

	var fileName string
	var offset int64
	for {
		n, err := fmt.Fscanf(file, "%s %d\n", &fileName, &offset)
		if n == 2 && err == nil {
			offsets[fileName] = offset
		} else if err != nil {
			break // Stop reading on error or EOF
		}
	}

	return offsets, nil
}

func Write(ctx context.Context, offsetChan chan OffsetData) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, stopping offset writing.")
			return
		case offsetData := <-offsetChan:
			if offsetData.Offset != 0 {
				offsetMap[offsetData.FileName] = offsetData.Offset // Update the offset map
				if time.Since(lastFlushTime) > 10*time.Second {
					err := _write_offset(offsetData)
					if err != nil {
						fmt.Printf("Error writing offset for %s: %v\n", offsetData.FileName, err)
						continue // Skip to the next iteration if there's an error
					}
					os.Rename(offsetFileTemp, offsetFilePath) // Rename the temp file to the final file
				}
			}
		}
	}
}

func _write_offset(offsetData OffsetData) error {
	// Open the file in append mode
	file, err := os.OpenFile(offsetFileTemp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		os.Create(offsetFileTemp) // Create the file if it doesn't exist
		file, err = os.OpenFile(offsetFileTemp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("error opening file %s: %w", offsetFileTemp, err)
		}
	}
	defer file.Close()

	for fileName, offset := range offsetMap {
		if _, err := file.WriteString(fmt.Sprintf("%s %d\n", fileName, offset)); err != nil {
			return fmt.Errorf("error writing to file %s: %w", offsetFileTemp, err)
		}
	}
	return nil
}

package offset

import (
	"fmt"
	"os"
	"test_gluent_mini/shared"
	"time"
)

const offsetFileTemp string = "./offset.tmp"
const offsetFilePath string = "./offset.state"

var lastFlushTime time.Time = time.Now()

func GetOffsetMap() (map[string]int64, error) {
	shared.M.Lock()
	defer shared.M.Unlock()
	if len(shared.OffsetMap) == 0 {
		offsets, err := _read()
		if err != nil {
			fmt.Printf("Error reading offsets: %v\n", err)
			return nil, err
		} else {
			for file, off := range offsets {
				shared.OffsetMap[file] = off
				fmt.Printf("Last offset for %s: %d\n", file, off)
			}
		}
	}
	return shared.OffsetMap, nil
}

func _read() (map[string]int64, error) {
	offsets := make(map[string]int64)
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

func Write() {
	for {
		select {
		case <-shared.Ctx.Done():
			fmt.Println("Context cancelled, stopping offset writing.")
			return
		case offsetData := <-shared.OffsetChannel:
			if offsetData.Offset != 0 {
				shared.M.Lock()
				shared.OffsetMap[offsetData.FileName] = offsetData.Offset // Update the offset map
				shared.M.Unlock()
				if time.Since(lastFlushTime) > 10*time.Second {
					err := _write_offset()
					if err != nil {
						fmt.Printf("Error writing offset for %s: %v\n", offsetData.FileName, err)
						continue // Skip to the next iteration if there's an error
					}
					if err := os.Rename(offsetFileTemp, offsetFilePath); err == nil {
						lastFlushTime = time.Now()
					} else {
						fmt.Printf("Error renaming temp offset file: %v\n", err)
					}
				}
			}
		}
	}
}

func _write_offset() error {
	shared.M.Lock()
	defer shared.M.Unlock()
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

	for fileName, offset := range shared.OffsetMap {
		if _, err := file.WriteString(fmt.Sprintf("%s %d\n", fileName, offset)); err != nil {
			return fmt.Errorf("error writing to file %s: %w", offsetFileTemp, err)
		}
	}
	fmt.Printf("Offsets written to %s successfully.\n", offsetFileTemp)
	return nil
}

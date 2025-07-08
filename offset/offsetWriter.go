package offset

import (
	"context"
	"fmt"
	"os"
)

var offsetFileTemp string = "./offset.tmp"
var offsetFilePath string = "./offset.state"

func WriterOffset(ctx context.Context, offsetChan chan int64) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, stopping offset writing.")
			return
		case offset := <-offsetChan:
			if offset != 0 {
				// Write the offset to the file
				if err := _write_offset(offset); err != nil {
					fmt.Printf("Error writing offset: %v\n", err)
					continue // Skip to the next iteration if there's an error
				}
				os.Rename(offsetFileTemp, offsetFilePath) // Rename the temp file to the final file
			}
		}
	}
}

func _write_offset(offset int64) error {
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

	// Write the offset to the file
	if _, err := file.WriteString(fmt.Sprintf("%d\n", offset)); err != nil {
		return fmt.Errorf("error writing to file %s: %w", offsetFileTemp, err)
	}

	return nil
}
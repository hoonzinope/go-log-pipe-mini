package offset

import (
	"fmt"
	"os"
)

func ReadOffset() (int64, error) {
	offsetFilePath := "./offset.state"
	file, err := os.Open(offsetFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil // If the file doesn't exist, return 0 offset
		}
		return 0, fmt.Errorf("error opening offset file %s: %w", offsetFilePath, err)
	}
	defer file.Close()

	var offset int64
	_, err = fmt.Fscanf(file, "%d", &offset)
	if err != nil {
		return 0, fmt.Errorf("error reading offset from file %s: %w", offsetFilePath, err)
	}

	return offset, nil
}
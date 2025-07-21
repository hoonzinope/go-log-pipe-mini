package input

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"test_gluent_mini/data"
	"test_gluent_mini/input/parse"
)

func TailFile(ctx context.Context,
	inputChan chan data.InputData,
	tag string, file string, parser string, offSetN int64) {
	newOffset := offSetN
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Tail the file
			inputDatas, lastOffset := _tail(tag, file, parser, newOffset)
			if lastOffset > newOffset {
				// Send the data to the input channel
				for _, inputData := range inputDatas {
					inputChan <- inputData
				}
				offsetChannel <- inputDatas[len(inputDatas)-1]
				newOffset = lastOffset // Update the offset
			}
		}
	}
}

func _tail(tag string, filePath string, parser string, offSetN int64) ([]data.InputData, int64) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return nil, offSetN
	}
	defer file.Close()

	// Move to the last known offset
	if _, err := file.Seek(offSetN, io.SeekStart); err != nil {
		fmt.Printf("Error seeking to offset %d in file %s: %v\n", offSetN, filePath, err)
		return nil, offSetN
	}

	scanner := bufio.NewScanner(file)
	var (
		results    []data.InputData
		lastOffset int64 = offSetN
	)
	for scanner.Scan() {
		line := scanner.Text()
		offset, _ := file.Seek(0, io.SeekCurrent)
		inputData := data.InputData{
			FileName: filePath,
			Tag:      tag,
			Raw:      line,
			Json:     nil,
			Offset:   offset,
		}
		if parser == "json" {
			inputData = _parseJSON(line, inputData)
		}
		results = append(results, inputData)
		lastOffset = offset
	}
	return results, lastOffset
}

func _parseJSON(line string, inputData data.InputData) data.InputData {
	if line == "" {
		return inputData // Return empty struct if input line is empty
	}
	jsonData, err := parse.ParseJSON(line)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return inputData
	}
	inputData.Json = jsonData
	return inputData
}

package input

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"test_gluent_mini/confmanager"
	"test_gluent_mini/data"
	"test_gluent_mini/input/parse"
	"test_gluent_mini/offset"
)

var offsetMap map[string]int64

// read offset of files
func init() {
	offsetMap = make(map[string]int64)
	offsets, err := offset.GetOffsetMap()
	if err != nil {
		fmt.Printf("Error reading offsets: %v\n", err)
	} else {
		for file, off := range offsets {
			offsetMap[file] = off
		}
	}
}

var inputFilePath string
var filePattern string
var tag string
var parser string

var cancel_ctx context.Context
var logLineChannel chan data.InputData
var offsetChannel chan data.InputData

func Configure(ctx context.Context,
	conf confmanager.Config,
	logLineChan chan data.InputData,
	offsetChan chan data.InputData) {
	inputFilePath = conf.Input.Path
	tag = conf.Input.Tag
	parser = conf.Input.Parser
	if inputFilePath == "" {
		fmt.Println("File path is not configured. Please check your configuration.")
		os.Exit(1)
	}
	filePattern = _dirToFilePattern(inputFilePath)

	cancel_ctx = ctx
	logLineChannel = logLineChan
	offsetChannel = offsetChan
}

func _dirToFilePattern(dir string) string {
	if strings.HasSuffix(dir, "/") || strings.HasSuffix(dir, "\\") {
		return dir + "*"
	}
	return dir + "/*"
}

var cancelMap = make(map[string]context.CancelFunc)

func ManagingNode() {
	files, err := filepath.Glob(filePattern)
	if err != nil {
		panic(fmt.Sprintf("Error reading files in path %s: %v", filePattern, err))
	}

	// add tail function of input files <- child process
	for _, file := range files {
		stat, err := os.Stat(file)
		if err != nil || !stat.Mode().IsRegular() {
			continue // 디렉터리, 링크 등은 건너뛰기
		}
		if _, exists := offsetMap[file]; !exists {
			offsetMap[file] = 0 // Initialize offset for new files
		}
		fileCtx, cancel := context.WithCancel(cancel_ctx)
		cancelMap[file] = cancel // Store the cancel function for later use
		go _tail(fileCtx, file, offsetMap[file])
	}
	// watch file changes and handle(add/remove) tail function dynamically
	go _watch(cancel_ctx, files)
}

func _tail(fileCtx context.Context, filePath string, offsetN int64) {
	for {
		select {
		case <-fileCtx.Done():
			fmt.Println("Context cancelled, stopping input reading.")
			return
		default:
			lastInputData := _tailFile(filePath, offsetN)
			if lastInputData.Raw != "" {
				logLineChannel <- lastInputData
				offsetChannel <- lastInputData
				offsetN = lastInputData.Offset // Update the offset for the next iteration
			}
		}
	}
}

func _tailFile(filePath string, offset int64) data.InputData {
	var inputData data.InputData = data.InputData{
		FileName: filePath,
		Tag:      tag,
		Raw:      "",
		Json:     nil,
		Offset:   offset,
	}
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return inputData
	}
	defer file.Close()

	// Move to the last known offset
	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		fmt.Printf("Error seeking to offset %d in file %s: %v\n", offset, filePath, err)
		return inputData
	}

	scanner := bufio.NewScanner(file)
	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return inputData
	}

	newOffset, _ := file.Seek(0, io.SeekCurrent) // Get the new offset after reading

	inputData.Raw = lastLine
	inputData.Offset = newOffset
	if parser == "json" {
		inputData = parseJSON(lastLine, inputData)
	}
	return inputData
}

func _watch(ctx context.Context, files []string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, stopping file watching.")
			return
		default:
			newFiles, err := _watchFiles(files)
			if err != nil {
				fmt.Printf("Error watching files: %v\n", err)
				return
			}
			files = newFiles
		}
	}
}

func _watchFiles(files []string) (newFiles []string, err error) {
	newFiles, err = filepath.Glob(filePattern)
	if err != nil {
		fmt.Printf("Error watching files in path %s: %v\n", filePattern, err)
		return nil, err
	}
	for _, file := range newFiles {
		if _, exists := offsetMap[file]; !exists {
			offsetMap[file] = 0 // Initialize offset for new files
			fileCtx, cancel := context.WithCancel(cancel_ctx)
			cancelMap[file] = cancel // Store the cancel function for later use
			go _tail(fileCtx, file, offsetMap[file])
		}
	}

	for _, file := range files {
		if _, exists := offsetMap[file]; !exists {
			continue
		}
		if !slices.Contains(newFiles, file) {
			cancel, exists := cancelMap[file]
			if exists {
				cancel()
				delete(cancelMap, file)
				delete(offsetMap, file)
				fmt.Printf("Stopped tailing file %s\n", file)
			}
		}
	}

	return newFiles, nil
}

func parseJSON(inputLine string, inputDataStruct data.InputData) data.InputData {
	if inputLine == "" {
		return inputDataStruct // Return empty struct if input line is empty
	}
	jsonData, err := parse.ParseJSON(inputLine)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return inputDataStruct
	}
	inputDataStruct.Json = jsonData
	return inputDataStruct
}

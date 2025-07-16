package input

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"test_gluent_mini/confmanager"
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
			fmt.Printf("Last offset for %s: %d\n", file, off)
		}
	}
}

var inputFilePath string
var name string
var cancel_ctx context.Context
var logLineChannel chan string
var offsetChannel chan offset.OffsetData

func Configure(ctx context.Context,
	conf confmanager.Config,
	logLineChan chan string,
	offsetChan chan offset.OffsetData) {
	inputFilePath = conf.Input.Path
	name = conf.Input.Name
	if inputFilePath == "" {
		fmt.Println("File path is not configured. Please check your configuration.")
		os.Exit(1)
	}
	cancel_ctx = ctx
	logLineChannel = logLineChan
	offsetChannel = offsetChan
}

var cancelMap = make(map[string]context.CancelFunc)

func ManagingNode() {
	if stat, _ := os.Stat(inputFilePath); stat.IsDir() {
		inputFilePath = filepath.Join(inputFilePath, "*")
	}
	files, err := filepath.Glob(inputFilePath)
	if err != nil {
		panic(fmt.Sprintf("Error reading files in path %s: %v", inputFilePath, err))
	}
	// add tail function of input files <- child process
	for _, file := range files {
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
			lastLine, newOffset := _tailFile(filePath, offsetN)
			lastLine = "[" + name + "] " + lastLine // Add name prefix to the log line
			if lastLine != "" {
				logLineChannel <- lastLine
				offsetData := offset.OffsetData{
					FileName: filePath,
					Offset:   newOffset,
				}
				offsetChannel <- offsetData
				offsetN = newOffset // Update the offset for the next iteration
			}
		}
	}
}

func _tailFile(filePath string, offset int64) (string, int64) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return "", offset
	}
	defer file.Close()

	// Move to the last known offset
	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		fmt.Printf("Error seeking to offset %d in file %s: %v\n", offset, filePath, err)
		return "", offset
	}

	scanner := bufio.NewScanner(file)
	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return "", offset
	}

	newOffset, _ := file.Seek(0, io.SeekCurrent) // Get the new offset after reading
	return lastLine, newOffset
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
	newFiles, err = filepath.Glob(inputFilePath)
	if err != nil {
		fmt.Printf("Error watching files in path %s: %v\n", inputFilePath, err)
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

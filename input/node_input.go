package input

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"test_gluent_mini/confmanager"
	"test_gluent_mini/shared"
)

func ManagingFileNode(
	inputConfig confmanager.InputConfig,
	inputChan chan shared.InputData) {

	conf := inputConfig

	_name := conf.Name
	_path := conf.Path
	_parser := conf.Parser

	filePatten := _dirToFilePattern(_path)
	files, err := filepath.Glob(filePatten)

	if err != nil {
		panic(fmt.Sprintf("Error finding files with pattern %s: %v", filePatten, err))
	}
	for _, file := range files {
		if !_isFile(file) {
			continue
		}
		shared.M.Lock()
		if _, exists := shared.OffsetMap[file]; !exists {
			shared.OffsetMap[file] = 0 // Initialize offset if not present
		}
		shared.M.Unlock()
		fileCtx, cancel := context.WithCancel(shared.Ctx)
		shared.CancelMap[file] = cancel
		// Start a goroutine to tail the file
		go TailFile(fileCtx, inputChan,
			_name, file, _parser, shared.OffsetMap[file])
	}
	//go _watch(cancelCtx, filePatten,inputChan, _name, _parser)
}

func _dirToFilePattern(dir string) string {
	if stat, err := os.Stat(dir); err == nil && stat.IsDir() {
		if !strings.Contains(dir, "*?") {
			return filepath.Join(dir, "*") // Assuming log files have .log extension
		}
	}
	return dir
}

func _isFile(path string) bool {
	stat, err := os.Stat(path)
	if err != nil || stat.IsDir() || !stat.Mode().IsRegular() {
		return false
	}
	return true
}

func _watch(ctx context.Context,
	filepattern string, inputChan chan shared.InputData,
	name string, parser string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Watch the files for changes
			_, err := _watchFiles(
				filepattern, inputChan, name, parser)
			if err != nil {
				fmt.Printf("Error watching files: %v\n", err)
				return
			}
		}
	}
}

func _watchFiles(filepattern string,
	inputChan chan shared.InputData,
	name string, parser string) ([]string, error) {

	shared.M.Lock()
	defer shared.M.Unlock()

	newFiles, err := filepath.Glob(filepattern)
	if err != nil {
		return nil, fmt.Errorf("error watching files: %v", err)
	}
	// new file go routine
	for _, file := range newFiles {
		if _, exists := shared.OffsetMap[file]; !exists {
			shared.OffsetMap[file] = 0 // Initialize offset for new files
			fileCtx, cancel := context.WithCancel(shared.Ctx)
			shared.CancelMap[file] = cancel // Store the cancel function for later use
			go TailFile(fileCtx,
				inputChan,
				name, file, parser, shared.OffsetMap[file])
		}
	}

	// delete not existing files
	var toDeleteFiles []string
	for file := range shared.OffsetMap {
		if !slices.Contains(newFiles, file) {
			toDeleteFiles = append(toDeleteFiles, file)
		}
	}
	if len(toDeleteFiles) != 0 {
		for _, file := range toDeleteFiles {
			cancel, exists := shared.CancelMap[file]
			if exists {
				cancel()                       // Cancel the context for the file
				delete(shared.CancelMap, file) // Remove from cancel map
				delete(shared.OffsetMap, file) // Remove from offset map
			}
		}
	}

	return newFiles, nil
}

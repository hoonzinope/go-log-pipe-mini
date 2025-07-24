package output

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"test_gluent_mini/shared"
	"time"
)

type FileOutput struct {
	Targets  []string
	Path     string
	Filename string
	Rolling  string
	MaxSize  string
	MaxFiles int
}

func (f FileOutput) _writeToFile(logLine shared.InputData) error {
	// Path/Filename 조합
	// rolling: default daily -> output.2023-10-01.log
	// max size: default 100MB -> 넘을경우 output.2023-10-01.1.log
	// max files: default 7 -> 넘을경우 oldest file 삭제

	// TODO : implement logic
	// output.log 파일을 열고, 해당 파일에 로그를 기록
	// 만약 rolling이 daily로 설정되어 있다면, 매일 새로운 파일을 생성
	// 만약 max size가 100MB로 설정되어 있다면, 파일 크기가 100MB를 초과할 경우 새로운 파일을 생성
	// 만약 max files가 7로 설정되어 있다면, path의 파일 개수 확인 및 초과시 가장 오래된 파일을 삭제

	var debug bool = true // 디버그 모드 설정

	var maxSize int64 = 100 * 1024 * 1024 // 100MB
	if strings.HasSuffix(f.MaxSize, "MB") {
		f.MaxSize = f.MaxSize[:len(f.MaxSize)-2]
		num, err := strconv.ParseInt(f.MaxSize, 10, 64) // Remove "MB"
		if err != nil {
			return fmt.Errorf("error parsing MaxSize %s: %v", f.MaxSize, err)
		}
		maxSize = num * 1024 * 1024 // 100MB
	}
	if strings.HasSuffix(f.MaxSize, "KB") {
		f.MaxSize = f.MaxSize[:len(f.MaxSize)-2] // Remove "KB"
		num, err := strconv.ParseInt(f.MaxSize, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing MaxSize %s: %v", f.MaxSize, err)
		}
		maxSize = num * 1024 // Convert to bytes
	}
	if strings.HasSuffix(f.MaxSize, "GB") {
		f.MaxSize = f.MaxSize[:len(f.MaxSize)-2] // Remove "GB"
		num, err := strconv.ParseInt(f.MaxSize, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing MaxSize %s: %v", f.MaxSize, err)
		}
		maxSize = num * 1024 * 1024 * 1024 // Convert to bytes
	}

	if _, err := os.ReadDir(f.Path); err != nil {
		if err := os.MkdirAll(f.Path, 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %v", f.Path, err)
		}
	}

	pattern := filepath.Join(f.Path, f.Filename, "*")
	if fileList, err := filepath.Glob(pattern); err == nil {
		if len(fileList) >= f.MaxFiles {
			// 파일 개수가 최대 개수를 초과한 경우, 가장 오래된 파일부터 삭제
			for _, file := range fileList[0 : len(fileList)-f.MaxFiles+1] {
				if err := os.Remove(file); err != nil {
					return fmt.Errorf("error removing old file %s: %v", file, err)
				}
			}
		}
	}

	file, err := os.OpenFile(filepath.Join(f.Path, f.Filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", f.Path+f.Filename, err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info for %s: %v", filepath.Join(f.Path, f.Filename), err)
	}
	switch f.Rolling {
	case "daily":
		currentTime := time.Now().Format("2006-01-02")
		if fileInfo.ModTime().Format("2006-01-02") != currentTime {
			file.Close()
			newFileName := fmt.Sprintf("%s.log.%s", f.Filename, currentTime)
			if err := os.Rename(filepath.Join(f.Path, f.Filename), filepath.Join(f.Path, newFileName)); err != nil {
				return fmt.Errorf("error renaming file: %v", err)
			}
			file, err = os.OpenFile(filepath.Join(f.Path, f.Filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("error opening new file %s: %v", filepath.Join(f.Path, f.Filename), err)
			}
		}
	case "hourly":
		currentTime := time.Now().Format("2006-01-02-15")
		if fileInfo.ModTime().Format("2006-01-02-15") != currentTime {
			file.Close()
			newFileName := fmt.Sprintf("%s.log.%s", f.Filename, currentTime)
			if err := os.Rename(filepath.Join(f.Path, f.Filename), filepath.Join(f.Path, newFileName)); err != nil {
				return fmt.Errorf("error renaming file: %v", err)
			}
			file, err = os.OpenFile(filepath.Join(f.Path, f.Filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("error opening new file %s: %v", filepath.Join(f.Path, f.Filename), err)
			}
		}
	case "monthly":
		currentTime := time.Now().Format("2006-01")
		if fileInfo.ModTime().Format("2006-01") != currentTime {
			file.Close()
			newFileName := fmt.Sprintf("%s.log.%s", f.Filename, currentTime)
			if err := os.Rename(filepath.Join(f.Path, f.Filename), filepath.Join(f.Path, newFileName)); err != nil {
				return fmt.Errorf("error renaming file: %v", err)
			}
			file, err = os.OpenFile(filepath.Join(f.Path, f.Filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("error opening new file %s: %v", filepath.Join(f.Path, f.Filename), err)
			}
		}
	}

	fileInfo, err = file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info for %s: %v", filepath.Join(f.Path, f.Filename), err)
	}
	if fileInfo.Size() >= maxSize { // 100MB
		file.Close()
		newFileName := fmt.Sprintf("%s.log.%d", f.Filename, time.Now().Unix())
		if err := os.Rename(filepath.Join(f.Path, f.Filename), filepath.Join(f.Path, newFileName)); err != nil {
			return fmt.Errorf("error renaming file: %v", err)
		}
		file, err = os.OpenFile(filepath.Join(f.Path, f.Filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("error opening new file %s: %v", filepath.Join(f.Path, f.Filename), err)
		}
	}

	if logLine.Json != nil {
		if _, err := file.WriteString(fmt.Sprintf("%s %s: %v\n", logLine.Tag, logLine.FileName, logLine.Json)); err != nil {
			file.Close()
			return fmt.Errorf("error writing JSON to file %s: %v", filepath.Join(f.Path, f.Filename), err)
		}
	} else if logLine.Raw != "" {
		if _, err := file.WriteString(fmt.Sprintf("%s %s: %s\n", logLine.Tag, logLine.FileName, logLine.Raw)); err != nil {
			file.Close()
			return fmt.Errorf("error writing raw log to file %s: %v", filepath.Join(f.Path, f.Filename), err)
		}
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("error closing file %s: %v", filepath.Join(f.Path, f.Filename), err)
	}
	if debug {
		fmt.Printf("Log written to file: %s\n", filepath.Join(f.Path, f.Filename))
	}
	return nil
}

func (f FileOutput) Out(ctx context.Context) {

	for _, target := range f.Targets {
		lineChan := shared.FilterChannel[target]
		go func(ctx context.Context, lineChan chan shared.InputData) {
			for {
				select {
				case <-ctx.Done():
					return
				case logLine := <-lineChan:
					if err := f._writeToFile(logLine); err != nil {
						fmt.Printf("Error writing to file: %v\n", err)
					}
					shared.OffsetChannel <- logLine
				}
			}
		}(ctx, lineChan)
	}
}

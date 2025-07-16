package generate

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func GenLogWithFolder(ctx context.Context) {
	var log_folder = "./logs"
	if _, err := os.Stat(log_folder); os.IsNotExist(err) {
		if err := os.Mkdir(log_folder, 0755); err != nil {
			fmt.Printf("Error creating log folder: %v\n", err)
			return
		}
	}
	var log_file_path string = fmt.Sprintf("%s/testlog.log", log_folder)
	var log_file_path1 string = fmt.Sprintf("%s/testlog1.log", log_folder)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, stopping log generation.")
			return
		default:
			_generate_log(log_file_path)
			_generate_log(log_file_path1)
			time.Sleep(1 * time.Second) // Sleep for 1 second before generating the next
		}
	}
}

var log_level = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

func _generate_log(filePath string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		os.Create(filePath)
	}
	defer file.Close()

	file.Write([]byte(_randomLogLine()))
}

func _randomLogLevel() string {
	return log_level[rand.Intn(len(log_level))]
}

func _randomLogLine() string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logLevel := _randomLogLevel()
	message := stringWithCharset(20, charset) // Generate a random message of 20 characters
	return fmt.Sprintf("%s [%s] %s\n", timestamp, logLevel, message)
}

// charset use random string
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// stringWithCharset return of random string
func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

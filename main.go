package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"test_gluent_mini/generate"
	"test_gluent_mini/input"
	"test_gluent_mini/offset"
	"test_gluent_mini/output"
)

var log_line_channel = make(chan string, 1000)
var offset_channel = make(chan int64, 1000)
var ctx, cancel = context.WithCancel(context.Background())
var wg sync.WaitGroup

func main() {
	fmt.Println("Hello, World!")

	lastOffset, err := offset.ReadOffset()
	if err != nil {
		fmt.Printf("Error reading offset: %v\n", err)
		os.Exit(1)
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-signalChan
		fmt.Println("Received shutdown signal, cleaning up...")
		cancel() // Cancel the context to stop all goroutines
		fmt.Println("Cleanup complete. Exiting.")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		input.TailFile(ctx, log_line_channel, "./testlog.log", lastOffset, offset_channel)
		fmt.Println("TailFile goroutine finished.")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		generate.GenLog(ctx)
		fmt.Println("GenLog goroutine finished.")
	}()
	
	wg.Add(1)
	go func() {
		defer wg.Done()
		output.Stdout(ctx, log_line_channel)
		fmt.Println("Stdout goroutine finished.")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		offset.WriterOffset(ctx, offset_channel)
		fmt.Println("WriterOffset goroutine finished.")
	}()

	wg.Wait()
}
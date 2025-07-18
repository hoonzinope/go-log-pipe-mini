package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"test_gluent_mini/confmanager"
	"test_gluent_mini/data"
	"test_gluent_mini/filter"
	"test_gluent_mini/generate"
	"test_gluent_mini/input"
	"test_gluent_mini/offset"
	"test_gluent_mini/output"
)

var log_line_channel = make(chan data.InputData, 1000)
var filter_line_channel = make(chan string, 1000)
var offset_channel = make(chan data.InputData, 1000)
var ctx, cancel = context.WithCancel(context.Background())
var wg sync.WaitGroup

func main() {
	fmt.Println("Starting the Gluent Mini application...")

	config, err := confmanager.ReadConfig("config.yml")
	if err != nil {
		fmt.Printf("Error reading configuration: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Configuration loaded: %+v\n", config)

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
		generate.GenerateJsonLog(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		input.Configure(ctx, config, log_line_channel, offset_channel)
		input.ManagingNode()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		filter.Configure(config)
		filter.FilterLine(ctx, log_line_channel, filter_line_channel)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		output.Configure(config)
		output.Out(ctx, filter_line_channel)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		offset.Write(ctx, offset_channel)
	}()

	wg.Wait()
}

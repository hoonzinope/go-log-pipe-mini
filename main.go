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

// var log_line_channel = make(chan data.InputData, 1000)

// var filter_line_channel = make(chan string, 1000)
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

	inputChannel := make(map[string]chan data.InputData)
	filterChannel := make(map[string]chan string)
	for _, inputConfig := range config.Inputs {
		inputChannel[inputConfig.Name] = make(chan data.InputData, 1000)
		filterChannel[inputConfig.Name] = make(chan string, 1000)
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
		generate.GenLogWithFolder(ctx)
		generate.GenerateJsonLog(ctx)
	}()

	input.Configure(ctx, config, inputChannel, offset_channel)
	wg.Add(1)
	go func() {
		defer wg.Done()
		input.ManagingNode()
	}()

	filter.Configure(ctx, config, inputChannel, filterChannel)
	wg.Add(1)
	go func() {
		defer wg.Done()
		filter.FilterLines()
	}()

	output.Configure(ctx, config, filterChannel)
	wg.Add(1)
	go func() {
		defer wg.Done()
		output.Out()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		offset.Write(ctx, offset_channel)
	}()

	wg.Wait()
}

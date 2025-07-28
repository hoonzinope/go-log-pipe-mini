package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"test_gluent_mini/confmanager"
	"test_gluent_mini/filter"
	"test_gluent_mini/generate"
	"test_gluent_mini/input"
	"test_gluent_mini/offset"
	"test_gluent_mini/output"
	"test_gluent_mini/server"
	"test_gluent_mini/shared"
)

func offsetInitialization() {
	_, err := offset.GetOffsetMap()
	if err != nil {
		fmt.Printf("Error getting offset map: %v\n", err)
		shared.Error_count.Add(1)
		return
	}
	fmt.Printf("Offset map initialized with %d entries.\n", len(shared.OffsetMap))
}

func main() {
	fmt.Println("Starting the Gluent Mini application...")
	offsetInitialization() // Initialize offsets from the offset file
	inputChannel := shared.InputChannel
	filterChannel := shared.FilterChannel
	_, cancel := shared.Ctx, shared.Cancel

	config, err := confmanager.ReadConfig("config.yml")
	if err != nil {
		fmt.Printf("Error reading configuration: %v\n", err)
		shared.Error_count.Add(1)
		os.Exit(1)
	}
	fmt.Printf("Configuration loaded: %+v\n", config)

	for _, inputConfig := range config.Inputs {
		inputChannel[inputConfig.Name] = make(chan shared.InputData, 1000)
		filterChannel[inputConfig.Name] = make(chan shared.InputData, 1000)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-signalChan
		fmt.Println("Received shutdown signal, cleaning up...")
		cancel() // Cancel the context to stop all goroutines
		fmt.Println("Cleanup complete. Exiting.")
	}()
	// Start the generate process for test data generation
	shared.Wg.Add(1)
	go func() {
		defer shared.Wg.Done()
		generate.GenLogWithFolder()
		generate.GenerateJsonLog()
	}()
	// Start the server to handle incoming requests for test data
	shared.Wg.Add(1)
	go func() {
		defer shared.Wg.Done()
		server.Run()
	}()

	input.Configure(config)
	shared.Wg.Add(1)
	go func() {
		defer shared.Wg.Done()
		input.ManagingNode()
	}()

	filter.Configure(config)
	shared.Wg.Add(1)
	go func() {
		defer shared.Wg.Done()
		filter.FilterLines()
	}()

	output.Configure(config)
	shared.Wg.Add(1)
	go func() {
		defer shared.Wg.Done()
		output.Out()
	}()

	shared.Wg.Add(1)
	go func() {
		defer shared.Wg.Done()
		offset.Write()
	}()

	shared.Wg.Add(1)
	go func() {
		defer shared.Wg.Done()
		shared.PrintStats(shared.Ctx)
	}()

	<-shared.Ctx.Done() // Wait for the context to be cancelled before proceeding
	fmt.Println("Context cancelled, shutting down gracefully...")
	server.Shutdown(shared.Ctx) // Ensure the server is properly shut down
	shared.Wg.Wait()
}

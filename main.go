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
	"test_gluent_mini/shared"
)

func offsetInitialization() {
	_, err := offset.GetOffsetMap()
	if err != nil {
		fmt.Printf("Error getting offset map: %v\n", err)
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

	shared.Wg.Add(1)
	go func() {
		defer shared.Wg.Done()
		generate.GenLogWithFolder()
		generate.GenerateJsonLog()
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

	shared.Wg.Wait()
}

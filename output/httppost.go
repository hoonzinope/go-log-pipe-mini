package output

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"test_gluent_mini/shared"
	"time"
)

type HttpOutput struct {
	Type           string
	Targets        []string
	Url            string
	Method         string
	Headers        map[string]string
	Timeout        string
	BATCH_SIZE     int
	FLUSH_INTERVAL string
}

func (h HttpOutput) Out(ctx context.Context, outputChannel map[string]chan shared.InputData) {

	if h.BATCH_SIZE == 0 {
		h.BATCH_SIZE = BATCH_SIZE_DEFAULT
	}
	if h.FLUSH_INTERVAL == "" {
		h.FLUSH_INTERVAL = FLUSH_INTERVAL_DEFAULT
	}
	duration, err := time.ParseDuration(h.FLUSH_INTERVAL)
	if err != nil {
		fmt.Printf("Error parsing FLUSH_INTERVAL %s: %v\n", h.FLUSH_INTERVAL, err)
		return
	}

	if h.Timeout == "" {
		h.Timeout = "5s" // Default timeout if not specified
	}
	timeout, err := time.ParseDuration(h.Timeout)
	if err != nil {
		fmt.Printf("Error parsing timeout %s: %v\n", h.Timeout, err)
		return
	}

	if err := h._waitForEndpointReady(timeout); err != nil {
		fmt.Printf("Error waiting for endpoint to be ready: %v\n", err)
		return
	}

	for _, target := range h.Targets {
		lineChan := outputChannel[target]
		go func(ctx context.Context, lineChan chan shared.InputData) {
			for {
				batch := make([]shared.InputData, 0, h.BATCH_SIZE)
				timer := time.NewTimer(duration)
			BATCHLOOP:
				for {
					select {
					case <-ctx.Done():
						return
					case logLine := <-lineChan:
						batch = append(batch, logLine)
						if len(batch) >= h.BATCH_SIZE {
							break BATCHLOOP
						}
					case <-timer.C:
						break BATCHLOOP
					}
				}

				for _, logLine := range batch {
					fmt.Printf("Sending HTTP POST to %s for target %s\n", h.Url, logLine.Tag)
					if err := h._writeToHttp(logLine, timeout); err != nil {
						fmt.Printf("Error writing to HTTP: %v\n", err)
					}
				}
				batch = nil // Clear the batch for the next iteration
				timer.Stop()
			}
		}(ctx, lineChan)
	}
}

func (h HttpOutput) _waitForEndpointReady(timeout time.Duration) error {
	// TODO: Implement health check logic to wait for the HTTP endpoint to be ready
	req, err := http.NewRequest("HEAD", h.Url, nil)
	if err != nil {
		return fmt.Errorf("Error creating HTTP request: %v", err)
	}
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP endpoint not ready, status: %s", resp.Status)
	}
	return nil
}

func (h HttpOutput) _writeToHttp(logLine shared.InputData, timeout time.Duration) error {
	// TODO : Implement the logic to convert logLine to the appropriate format (e.g., JSON)

	// Example implementation (replace with actual logic)
	jsonData, err := json.Marshal(logLine)
	if err != nil {
		return fmt.Errorf("Error marshaling logLine to JSON: %v", err)
	}

	// Send the JSON data to the HTTP endpoint
	// add timeout to the request
	client := &http.Client{
		Timeout: timeout, // Set a timeout for the request
	}
	resp, err := client.Post(h.Url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Error sending HTTP POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status: %s", resp.Status)
	}
	fmt.Printf("HTTP POST successful for target %s: %s\n", logLine.Tag, logLine.FileName)
	return nil
}

package server

import (
	"fmt"
	"testing"
)

func TestRunServer(t *testing.T) {
	// Initialize the server
	fmt.Println("Running server tests...")
	debug := true // Set to true to enable debug mode
	Run(debug)
}

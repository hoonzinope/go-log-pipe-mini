package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

var httpServer *http.Server

func Run(debug bool) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set
	}
	httpServer = &http.Server{
		Addr: ":" + port,
	}
	httpServer.Handler = http.NewServeMux()
	mux := httpServer.Handler.(*http.ServeMux)
	healthCheck(mux)
	if debug {
		handleLogRequests(mux)
	}
	endpoint_metrics(mux)
	fmt.Println("Starting log server on port =>", port)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("Error starting log server:", err)
		return
	}
	fmt.Println("Log server started successfully")
}

func Shutdown(ctx context.Context) {
	fmt.Println("Shutting down log server...")
	if err := httpServer.Shutdown(ctx); err != nil {
		fmt.Println("Error shutting down log server:", err)
	}
}

package server

import (
	"context"
	"fmt"
	"net/http"
)

var httpServer *http.Server

func Run() {
	httpServer = &http.Server{
		Addr: ":8080",
	}
	httpServer.Handler = http.NewServeMux()
	mux := httpServer.Handler.(*http.ServeMux)
	healthCheck(mux)
	handleLogRequests(mux)
	endpoint_metrics(mux)
	fmt.Println("Starting log server on :8080")
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

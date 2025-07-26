package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

var httpServer *http.Server

func Run() {
	httpServer = &http.Server{
		Addr: ":8080",
	}
	handleRequests()
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

func handleRequests() {
	httpServer.Handler = http.NewServeMux()
	httpServer.Handler.(*http.ServeMux).HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the Log Server!")
	})
	httpServer.Handler.(*http.ServeMux).HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Log Server is healthy")
	})
	httpServer.Handler.(*http.ServeMux).HandleFunc("/logs", logHandler)
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		fmt.Println("Error reading request body:", err)
		return
	}
	fmt.Println("Received POST:", string(body))
	w.WriteHeader(http.StatusOK)
}

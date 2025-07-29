package server

import (
	"fmt"
	"io"
	"net/http"
)

func handleLogRequests(mux *http.ServeMux) {
	mux.HandleFunc("/logs", logHandler)
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

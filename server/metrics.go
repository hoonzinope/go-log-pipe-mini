package server

import (
	"fmt"
	"net/http"
	"test_gluent_mini/shared"
)

func endpoint_metrics(mux *http.ServeMux) {
	mux.HandleFunc("/metrics", _metrics)
}

func _metrics(w http.ResponseWriter, r *http.Request) {
	output := "Log Pipe Mini Metrics\n"
	output += "[stat] Input: " + fmt.Sprint(shared.Input_count.Load()) + "\n"
	output += "[stat] Filter: " + fmt.Sprint(shared.Filter_count.Load()) + "\n"
	output += "[stat] Output: " + fmt.Sprint(shared.Output_count.Load()) + "\n"
	output += "[stat] Error: " + fmt.Sprint(shared.Error_count.Load()) + "\n"
	output += "[stat] Avg Latency: " + shared.GetAverageLatency().String() + "\n"
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(output))
}

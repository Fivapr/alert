package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	gauge   map[string][]float64
	counter map[string]int64
}

// Global storage instance
var storage = &MemStorage{
	gauge:   make(map[string][]float64),
	counter: make(map[string]int64),
}

func updateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Split the URL path into components
	pathComponents := strings.Split(r.URL.Path, "/")

	// Check if the path has the expected number of components ("/update/<METRIC_TYPE>/<METRIC_NAME>/<METRIC_VALUE>")
	// Note: The first element after split will be empty due to leading "/", hence expecting 5 components
	if len(pathComponents) != 5 {
		http.Error(w, "Path should be /update/<METRIC_TYPE>/<METRIC_NAME>/<METRIC_VALUE>", http.StatusNotFound)
		return
	}

	metricType := pathComponents[2]
	metricName := pathComponents[3]
	metricValueStr := pathComponents[4]

	if metricType != "gauge" && metricType != "counter" {
		http.Error(w, "Path should be /update/<METRIC_TYPE>/<METRIC_NAME>/<METRIC_VALUE>", http.StatusBadRequest)
		return
	}

	if metricType == "gauge" {
		metricValue, err := strconv.ParseFloat(metricValueStr, 64)
		if err != nil {
			http.Error(w, "Invalid metric value", http.StatusBadRequest)
			return
		}

		storage.gauge[metricName] = append(storage.gauge[metricName], metricValue)
	}

	if metricType == "counter" {
		metricValue, err := strconv.ParseInt(metricValueStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid metric value", http.StatusBadRequest)
			return
		}

		_, ok := storage.counter[metricName]
		if ok {
			storage.counter[metricName] += metricValue
		} else {
			storage.counter[metricName] = metricValue
		}
	}

}

func main() {
	http.HandleFunc("/update/", updateMetric)

	fmt.Println("Server is listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

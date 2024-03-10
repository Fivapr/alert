package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
	"os"
	"strconv"
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

func getMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	w.Header().Set("Content-Type", "text/plain")

	switch metricType {
	case "gauge":
		if values, ok := storage.gauge[metricName]; ok {
			latestValue := values[len(values)-1]
			fmt.Fprintf(w, "%v", latestValue)
		} else {
			http.Error(w, "Metric not found", http.StatusNotFound)
		}

	case "counter":
		if value, ok := storage.counter[metricName]; ok {
			fmt.Fprintf(w, "%v", value)
		} else {
			http.Error(w, "Metric not found", http.StatusNotFound)
		}

	default:
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
	}
}

func updateMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValueStr := chi.URLParam(r, "metricValue")

	switch metricType {
	case "gauge":
		metricValue, err := strconv.ParseFloat(metricValueStr, 64)
		if err != nil {
			http.Error(w, "Invalid metric value", http.StatusBadRequest)
			return
		}
		storage.gauge[metricName] = append(storage.gauge[metricName], metricValue)

	case "counter":
		metricValue, err := strconv.ParseInt(metricValueStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid metric value", http.StatusBadRequest)
			return
		}
		storage.counter[metricName] += metricValue

	default:
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Updated %s metric %s with value %s\n", metricType, metricName, metricValueStr)
}

const metricsTemplateStr = `
<!DOCTYPE html>
<html>
<head>
    <title>Metrics</title>
</head>
<body>
    <h1>Metrics</h1>
    <h2>Gauges</h2>
    <ul>
        {{range $name, $values := .Gauge}}
        <li>{{$name}}: {{$values}}</li>
        {{end}}
    </ul>
    <h2>Counters</h2>
    <ul>
        {{range $name, $value := .Counter}}
        <li>{{$name}}: {{$value}}</li>
        {{end}}
    </ul>
</body>
</html>
`

func getAll(w http.ResponseWriter, r *http.Request) {
	// Parse the template
	tmpl, err := template.New("metrics").Parse(metricsTemplateStr)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	// Data to pass to the template
	data := struct {
		Gauge   map[string][]float64
		Counter map[string]int64
	}{
		Gauge:   storage.gauge,
		Counter: storage.counter,
	}

	// Execute the template, writing the generated HTML to the response
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

var aFlag = *flag.String("a", "localhost:8080", "Port to run the server on")

func main() {
	// Custom usage function to provide detailed help text
	// This does not change the exit code behavior but improves user guidance
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	defaultAddr := os.Getenv("ADDRESS")
	if defaultAddr != "" {
		aFlag = defaultAddr
	}

	r := chi.NewRouter()

	r.Get("/", getAll)
	r.Get("/value/{metricType}/{metricName}", getMetric)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", updateMetric)

	fmt.Printf("Server is listening on :%s\n", aFlag)
	if err := http.ListenAndServe(aFlag, r); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start server: %v\n", err)
		os.Exit(1)
	}
}

package main

import (
	"flag"
	"fmt"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

var snapshot MemStatsSnapshot
var pollCount int
var randomValue int

type MemStatsSnapshot struct {
	Alloc         uint64  `json:"Alloc"`
	BuckHashSys   uint64  `json:"BuckHashSys"`
	Frees         uint64  `json:"Frees"`
	GCCPUFraction float64 `json:"GCCPUFraction"`
	GCSys         uint64  `json:"GCSys"`
	HeapAlloc     uint64  `json:"HeapAlloc"`
	HeapIdle      uint64  `json:"HeapIdle"`
	HeapInuse     uint64  `json:"HeapInuse"`
	HeapObjects   uint64  `json:"HeapObjects"`
	HeapReleased  uint64  `json:"HeapReleased"`
	HeapSys       uint64  `json:"HeapSys"`
	LastGC        uint64  `json:"LastGC"`
	Lookups       uint64  `json:"Lookups"`
	MCacheInuse   uint64  `json:"MCacheInuse"`
	MCacheSys     uint64  `json:"MCacheSys"`
	MSpanInuse    uint64  `json:"MSpanInuse"`
	MSpanSys      uint64  `json:"MSpanSys"`
	Mallocs       uint64  `json:"Mallocs"`
	NextGC        uint64  `json:"NextGC"`
	NumForcedGC   uint32  `json:"NumForcedGC"`
	NumGC         uint32  `json:"NumGC"`
	OtherSys      uint64  `json:"OtherSys"`
	PauseTotalNs  uint64  `json:"PauseTotalNs"`
	StackInuse    uint64  `json:"StackInuse"`
	StackSys      uint64  `json:"StackSys"`
	Sys           uint64  `json:"Sys"`
	TotalAlloc    uint64  `json:"TotalAlloc"`
}

func createMemStatsSnapshot(m runtime.MemStats) MemStatsSnapshot {
	snapshot := MemStatsSnapshot{
		Alloc:         m.Alloc,
		BuckHashSys:   m.BuckHashSys,
		Frees:         m.Frees,
		GCCPUFraction: m.GCCPUFraction,
		GCSys:         m.GCSys,
		HeapAlloc:     m.HeapAlloc,
		HeapIdle:      m.HeapIdle,
		HeapInuse:     m.HeapInuse,
		HeapObjects:   m.HeapObjects,
		HeapReleased:  m.HeapReleased,
		HeapSys:       m.HeapSys,
		LastGC:        m.LastGC,
		Lookups:       m.Lookups,
		MCacheInuse:   m.MCacheInuse,
		MSpanInuse:    m.MSpanInuse,
		MSpanSys:      m.MSpanSys,
		Mallocs:       m.Mallocs,
		NextGC:        m.NextGC,
		NumForcedGC:   m.NumForcedGC,
		NumGC:         m.NumGC,
		OtherSys:      m.OtherSys,
		PauseTotalNs:  m.PauseTotalNs,
		StackInuse:    m.StackInuse,
		StackSys:      m.StackSys,
		Sys:           m.Sys,
		TotalAlloc:    m.TotalAlloc,
	}

	return snapshot
}

func sendMetric(metricType string, metricName string, metricValue string) {
	client := resty.New()

	url := fmt.Sprintf("http://%s/update/%s/%s/%s", aFlag, metricType, metricName, metricValue)

	resp, err := client.R().
		SetHeader("Content-Type", "text/plain").
		Post(url)

	if err != nil {
		fmt.Printf("Error sending request to %s: %v\n", url, err)
		return
	}

	fmt.Printf("Successfully sent metric %s to server. Status Code: %d\n", metricName, resp.StatusCode())
}

func sendMetrics() {
	v := reflect.ValueOf(snapshot)
	typeOfSnapshot := v.Type()

	for i := 0; i < v.NumField(); i++ {
		metricName := typeOfSnapshot.Field(i).Name
		metricValue := fmt.Sprintf("%v", v.Field(i).Interface())
		metricType := "gauge" // Assuming all metrics are gauges

		sendMetric(metricType, metricName, metricValue)
	}

	sendMetric("gauge", "RandomValue", strconv.Itoa(randomValue))
	sendMetric("counter", "PollCount", strconv.Itoa(pollCount))
}

func updateMemStatsPeriodically() {
	pollTicker := time.NewTicker(time.Duration(pFlag) * time.Second)

	go func() {
		for range pollTicker.C {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			snapshot = createMemStatsSnapshot(m)
			randomValue = rand.Intn(1000)
			pollCount++
		}
	}()
}

func sendMetricsPeriodically() {
	reportTicker := time.NewTicker(time.Duration(rFlag) * time.Second)

	go func() {
		for range reportTicker.C {
			sendMetrics()
		}
	}()
}

var (
	aFlag = *flag.String("a", "localhost:8080", "Port to run the server on")
	pFlag = *flag.Int("p", 2, "poll interval")
	rFlag = *flag.Int("r", 10, "report interval")
)

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

	defaultRepInterval := os.Getenv("REPORT_INTERVAL")
	if defaultAddr != "" {
		if repInt, err := strconv.Atoi(defaultRepInterval); err == nil {
			rFlag = repInt
		}
	}

	defaultPollInterval := os.Getenv("POLL_INTERVAL")
	if defaultAddr != "" {
		if pollInt, err := strconv.Atoi(defaultPollInterval); err == nil {
			pFlag = pollInt
		}
	}

	updateMemStatsPeriodically()
	sendMetricsPeriodically()
	select {}
}

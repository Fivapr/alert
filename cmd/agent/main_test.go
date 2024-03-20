package main

import (
	"reflect"
	"runtime"
	"testing"
)

func Test_createMemStatsSnapshot(t *testing.T) {
	type args struct {
		m runtime.MemStats
	}
	tests := []struct {
		name string
		args args
		want MemStatsSnapshot
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createMemStatsSnapshot(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createMemStatsSnapshot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sendMetric(t *testing.T) {
	type args struct {
		metricType  string
		metricName  string
		metricValue string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendMetric(tt.args.metricType, tt.args.metricName, tt.args.metricValue)
		})
	}
}

func Test_sendMetrics(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendMetrics()
		})
	}
}

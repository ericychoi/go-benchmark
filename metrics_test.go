package main

import (
	metrics "github.com/rcrowley/go-metrics"

	"strconv"
	"testing"
)

// % go test -test.bench=BenchmarkMetrics -run=XXX
// goos: darwin
// goarch: amd64
// BenchmarkMetricsGetOrRegister-8          1000000              1036 ns/op
//
// BenchmarkMetricsPreRegistered-8         10000000               216 ns/op
// PASS
// ok      _/Users/ericchoi/git/go-benchmark       9.448s

var reg metrics.Registry

func BenchmarkMetricsGetOrRegister(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gauge := metrics.GetOrRegisterGauge(strconv.Itoa(i), reg)
		gauge.Update(int64(i))
	}
}

func BenchmarkMetricsPreRegistered(b *testing.B) {
	pre := make(map[string]metrics.Gauge)
	for i := 0; i < b.N; i++ {
		pre[strconv.Itoa(i)] = metrics.NewGauge()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pre[strconv.Itoa(i)].Update(int64(i))
	}
}

package main

import (
	"strconv"
	"testing"
)

func BenchmarkMapDynamicInsert(b *testing.B) {
	// fill up the map
	structMap := make(map[string]struct{}, 100)
	for i := 0; i < 100; i++ {
		a := strconv.Itoa(i)
		structMap[a] = struct{}{}
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		a := strconv.Itoa(i + 100)
		structMap[a] = struct{}{}
	}
}

func BenchmarkMapStaticInsert(b *testing.B) {
	// fill up the map
	structMap := make(map[string]struct{}, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		a := strconv.Itoa(i)
		structMap[a] = struct{}{}
	}
}

func benchmarkMapInit(size int, b *testing.B) {
	var structMap map[string]struct{}
	for i := 0; i < b.N; i++ {
		structMap = make(map[string]struct{}, size)
	}
	structMap["a"] = struct{}{}
}

func BenchmarkMapInit10(b *testing.B)      { benchmarkMapInit(10, b) }
func BenchmarkMapInit100(b *testing.B)     { benchmarkMapInit(100, b) }
func BenchmarkMapInit1000(b *testing.B)    { benchmarkMapInit(1000, b) }
func BenchmarkMapInit10000(b *testing.B)   { benchmarkMapInit(10000, b) }
func BenchmarkMapInit100000(b *testing.B)  { benchmarkMapInit(100000, b) }
func BenchmarkMapInit1000000(b *testing.B) { benchmarkMapInit(1000000, b) }

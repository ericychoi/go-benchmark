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

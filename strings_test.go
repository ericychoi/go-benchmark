package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

//output
// % go test -bench=BenchmarkStrings
// testing: warning: no tests to run
// PASS
// BenchmarkStringsCasting-8	    3000	    492737 ns/op
// BenchmarkStringsBytes-8  	2000000000	         0.29 ns/op
// ok  	_/Users/ericchoi/git/go-benchmark	7.849s

func BenchmarkStringsCasting(b *testing.B) {
	var input []byte

	for i := 0; i < 1000000; i++ {
		a := strconv.Itoa(i)
		input = append(input, []byte(a)...)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dummyFunction2(string(input))
	}
}

func BenchmarkStringsBytes(b *testing.B) {
	var input []byte

	for i := 0; i < 1000000; i++ {
		a := strconv.Itoa(i)
		input = append(input, []byte(a)...)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dummyFunction1(input)
	}
}

func dummyFunction2(s string) {
}

func dummyFunction1(b []byte) {
}

// stringBuilder vs concatenation when number of pieces is fixed
var numPieces = 4

func BenchmarkStringsBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var b strings.Builder
		for i := 0; i <= numPieces; i++ {
			// use fmt.Printf to simulate the same condition as below)
			fmt.Fprint(&b, fmt.Sprintf("%d", i))
		}
		fmt.Fprint(ioutil.Discard, b.String())
	}
}

func BenchmarkStringsConcatenate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var str string
		for i := 0; i <= numPieces; i++ {
			str = str + fmt.Sprintf("%d", i)
		}
		fmt.Fprint(ioutil.Discard, str)
	}
}

// Bitflags vs string comparison for type checking
const (
	StatusPending   = 1 << 0
	StatusProcessed = 1 << 1
	StatusFailed    = 1 << 2
	StatusCompleted = 1 << 3
)

func BenchmarkStringComparison(b *testing.B) {
	statuses := []string{"pending", "processed", "failed", "completed", "pending", "processed"}
	var count int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		status := statuses[i%len(statuses)]
		if status == "processed" || status == "completed" {
			count++
		}
	}
	_ = count
}

func BenchmarkBitflagComparison(b *testing.B) {
	statuses := []uint8{StatusPending, StatusProcessed, StatusFailed, StatusCompleted, StatusPending, StatusProcessed}
	var count int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		status := statuses[i%len(statuses)]
		if status&(StatusProcessed|StatusCompleted) != 0 {
			count++
		}
	}
	_ = count
}

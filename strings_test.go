package main

import (
	"strconv"
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

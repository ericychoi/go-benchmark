package main

import (
	"regexp"
	"strings"
	"testing"
)

// during a set comparison, would a regex ("a|b|c") be faster than a loop?
// turns out a loop is still faster, assuming a small set of needles

// goos: darwin
// goarch: amd64
// pkg: github.com/ericychoi/go-benchmark
// cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
// BenchmarkSetComparisonLoop-8   	  543616	      2135 ns/op	    1009 B/op	       9 allocs/op
func BenchmarkSetComparisonLoop(b *testing.B) {
	stack := []string{"abc", "def", "ghi", "jkl", "mno", "pqr", "stu", "vwx", "yz"} // stack
	needle := []string{"a", "d", "g"}                                               // needle
	for i := 0; i < b.N; i++ {
		count := 0
		for _, s := range stack {
			for _, t := range needle {
				if strings.Contains(s, t) {
					count++
				}
			}
		}
		if count == 3 {
			b.Logf("found all")
		}
	}
}

// goos: darwin
// goarch: amd64
// pkg: github.com/ericychoi/go-benchmark
// cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
// BenchmarkSetComparisonRegex-8   	  414184	      2834 ns/op	    1021 B/op	       9 allocs/op
func BenchmarkSetComparisonRegex(b *testing.B) {
	stack := []string{"abc", "def", "ghi", "jkl", "mno", "pqr", "stu", "vwx", "yz"} // stack
	needle := []string{"a", "d", "g"}                                               // needle
	needleRegex := "(?i)" + strings.Join(needle, "|")
	regex := regexp.MustCompile(needleRegex)
	for i := 0; i < b.N; i++ {
		count := 0
		for _, s := range stack {
			if regex.MatchString(s) {
				count++
			}
		}
		if count == 3 {
			b.Logf("found all")
		}
	}
}

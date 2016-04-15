package main

import "testing"

//output
// % go test -bench=BenchmarkVar
// testing: warning: no tests to run
// PASS
// BenchmarkVarLoopDecl-8	2000000000	         0.54 ns/op
// BenchmarkVarDecl-8    	2000000000	         0.53 ns/op
// ok  	_/Users/ericchoi/git/go-benchmark	5.268s

func BenchmarkVarLoopDecl(b *testing.B) {
	for i := 0; i < b.N; i++ {
		a := i + 1
		a++ // so that compiler won't complain about unused a
	}
}

func BenchmarkVarDecl(b *testing.B) {
	var a int
	for i := 0; i < b.N; i++ {
		a = i + 1
		a++ // so that compiler won't complain about unused a
	}
}

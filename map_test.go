package main

import (
	"math/rand"
	"os"
	"strconv"
	"testing"
)

var isVerbose bool
var stringMap map[string]string
var structMap map[KeyStruct]string

type KeyStruct struct {
	A string
	B string
}

func TestMain(m *testing.M) {
	// encode given file into base64
	stringMap = make(map[string]string, 100)
	structMap = make(map[KeyStruct]string, 100)

	for i := 0; i < 50; i++ {
		a, b := strconv.Itoa(i), strconv.Itoa(i*2)
		structMap[KeyStruct{a, b}] = "ddd"
	}

	for i := 0; i < 50; i++ {
		a, b := strconv.Itoa(i), strconv.Itoa(i*2)
		stringMap[a+"_"+b] = "ddd"
	}

	os.Exit(m.Run())
}

func BenchmarkAnonymousMapLookup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		i, j := strconv.Itoa(rand.Intn(50)), strconv.Itoa(rand.Intn(50))
		_, ok := structMap[KeyStruct{i, j}]
		if !ok {
			b.Logf("cannot retrieve value in structMap %s, %s\n", i, j)
		}
	}
}

func BenchmarkStringLookup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		i, j := strconv.Itoa(rand.Intn(50)), strconv.Itoa(rand.Intn(50))
		_, ok := stringMap[i+"_"+j]
		if !ok {
			b.Logf("cannot retrieve value in stringMap %s, %s\n", i, j)
		}
	}
}

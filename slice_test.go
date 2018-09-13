package main

// % go test -bench=BenchmarkErrorSlice -benchmem
// goos: darwin
// goarch: amd64
// BenchmarkErrorSlice-8   	 3000000	       457 ns/op	      86 B/op	       6 allocs/op
// PASS
// ok  	_/Users/ericchoi/git/go-benchmark	1.956s

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
)

type Mock struct {
	client MockClient
}

type MockClient struct{}

func (c *MockClient) Quit() error {
	if rand.Intn(10) < 5 {
		return errors.New("err")
	}
	return nil
}

func (c *MockClient) Close() error {
	if rand.Intn(10) < 5 {
		return errors.New("err")
	}
	return nil
}

func (c *Mock) ProblemSnippet() error {
	var err error
	errors := []error{}

	err = c.client.Quit()
	if err != nil {
		errors = append(errors, err)
	}

	// even if quit fails, we still want to try to close the connection
	err = c.client.Close()
	if err != nil {
		errors = append(errors, err)
	}

	errStr := ""
	for _, e := range errors {
		errStr += fmt.Sprintf("%s\n", e.Error())
	}

	if errStr != "" {
		return fmt.Errorf("%s", errStr)
	}

	return nil
}

func BenchmarkErrorSlice(b *testing.B) {
	c := &Mock{}
	for i := 0; i < b.N; i++ {
		c.ProblemSnippet()
	}
}

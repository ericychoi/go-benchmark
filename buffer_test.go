package main

import (
	"io"
	"os"
	"sync"
	"testing"
	"time"
)

// slowReader wraps an io.Reader and adds a delay before each read
type slowReader struct {
	r     io.Reader
	delay time.Duration
}

func (s *slowReader) Read(p []byte) (n int, err error) {
	time.Sleep(s.delay)
	return s.r.Read(p)
}

// BenchmarkBufferPreallocated reads data into a pre-allocated 1KB buffer
func BenchmarkBufferPreallocated(b *testing.B) {
	f, err := os.Open("/dev/random")
	if err != nil {
		b.Fatalf("failed to open /dev/random: %v", err)
	}
	defer f.Close()

	const bufSize = 1024
	const readSize = 1024 // read 1KB per iteration

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, bufSize)
		n, err := io.ReadFull(f, buf)
		if err != nil {
			b.Fatalf("read failed: %v", err)
		}
		if n != readSize {
			b.Fatalf("expected %d bytes, got %d", readSize, n)
		}
	}
}

// BenchmarkBufferDynamicGrow reads data into a dynamically growing buffer starting at 10 bytes
func BenchmarkBufferDynamicGrow(b *testing.B) {
	f, err := os.Open("/dev/random")
	if err != nil {
		b.Fatalf("failed to open /dev/random: %v", err)
	}
	defer f.Close()

	const readSize = 1024 // target 1KB per iteration

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, 10) // start small
		totalRead := 0

		for totalRead < readSize {
			// Grow buffer if needed
			if len(buf) < readSize {
				newBuf := make([]byte, len(buf)*2)
				copy(newBuf, buf[:totalRead])
				buf = newBuf
			}

			n, err := f.Read(buf[totalRead:])
			if err != nil {
				b.Fatalf("read failed: %v", err)
			}
			totalRead += n
		}
	}
}

// BenchmarkBufferDynamicGrowOptimized uses append pattern for dynamic growth
func BenchmarkBufferDynamicGrowOptimized(b *testing.B) {
	f, err := os.Open("/dev/random")
	if err != nil {
		b.Fatalf("failed to open /dev/random: %v", err)
	}
	defer f.Close()

	const readSize = 1024 // target 1KB per iteration
	const chunkSize = 10  // read in small chunks

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, 0, 10) // start small with capacity 10
		chunk := make([]byte, chunkSize)

		for len(buf) < readSize {
			n, err := f.Read(chunk)
			if err != nil {
				b.Fatalf("read failed: %v", err)
			}
			buf = append(buf, chunk[:n]...)
		}
	}
}

// BenchmarkBufferPreallocatedSlow reads with 1ms delay per read
func BenchmarkBufferPreallocatedSlow(b *testing.B) {
	f, err := os.Open("/dev/random")
	if err != nil {
		b.Fatalf("failed to open /dev/random: %v", err)
	}
	defer f.Close()

	slow := &slowReader{r: f, delay: time.Millisecond}

	const bufSize = 1024
	const readSize = 1024

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, bufSize)
		n, err := io.ReadFull(slow, buf)
		if err != nil {
			b.Fatalf("read failed: %v", err)
		}
		if n != readSize {
			b.Fatalf("expected %d bytes, got %d", readSize, n)
		}
	}
}

// BenchmarkBufferDynamicGrowSlow reads with 1ms delay per read
func BenchmarkBufferDynamicGrowSlow(b *testing.B) {
	f, err := os.Open("/dev/random")
	if err != nil {
		b.Fatalf("failed to open /dev/random: %v", err)
	}
	defer f.Close()

	slow := &slowReader{r: f, delay: time.Millisecond}

	const readSize = 1024

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, 10)
		totalRead := 0

		for totalRead < readSize {
			if len(buf) < readSize {
				newBuf := make([]byte, len(buf)*2)
				copy(newBuf, buf[:totalRead])
				buf = newBuf
			}

			n, err := slow.Read(buf[totalRead:])
			if err != nil {
				b.Fatalf("read failed: %v", err)
			}
			totalRead += n
		}
	}
}

// BenchmarkBufferDynamicGrowOptimizedSlow uses append pattern with 1ms delay
func BenchmarkBufferDynamicGrowOptimizedSlow(b *testing.B) {
	f, err := os.Open("/dev/random")
	if err != nil {
		b.Fatalf("failed to open /dev/random: %v", err)
	}
	defer f.Close()

	slow := &slowReader{r: f, delay: time.Millisecond}

	const readSize = 1024
	const chunkSize = 10

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, 0, 10)
		chunk := make([]byte, chunkSize)

		for len(buf) < readSize {
			n, err := slow.Read(chunk)
			if err != nil {
				b.Fatalf("read failed: %v", err)
			}
			buf = append(buf, chunk[:n]...)
		}
	}
}

// BenchmarkBufferParallelReadSlow reads small chunks in parallel with 1ms delay
func BenchmarkBufferParallelReadSlow(b *testing.B) {
	f, err := os.Open("/dev/random")
	if err != nil {
		b.Fatalf("failed to open /dev/random: %v", err)
	}
	defer f.Close()

	slow := &slowReader{r: f, delay: time.Millisecond}

	const readSize = 1024
	const chunkSize = 64
	const numChunks = readSize / chunkSize

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		type result struct {
			idx  int
			data []byte
			err  error
		}

		results := make(chan result, numChunks)
		var wg sync.WaitGroup

		// Launch parallel reads
		for j := 0; j < numChunks; j++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				chunk := make([]byte, chunkSize)
				n, err := io.ReadFull(slow, chunk)
				if err != nil {
					results <- result{idx: idx, err: err}
					return
				}
				results <- result{idx: idx, data: chunk[:n]}
			}(j)
		}

		wg.Wait()
		close(results)

		// Assemble buffer
		buf := make([]byte, 0, readSize)
		chunks := make([][]byte, numChunks)
		for r := range results {
			if r.err != nil {
				b.Fatalf("read failed: %v", r.err)
			}
			chunks[r.idx] = r.data
		}

		for _, chunk := range chunks {
			buf = append(buf, chunk...)
		}
	}
}

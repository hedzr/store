package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(1)

	start := time.Now()

	go calc()

	calc()

	elapsed := time.Since(start)
	_, _ = fmt.Printf("Time elapsed: %v | switch cost ~ %vns\n", elapsed, elapsed.Nanoseconds()/Counter)
}

func calc() {
	for i := 0; i < Counter; i++ {
		runtime.Gosched()
	}
}

const Counter = 10000000

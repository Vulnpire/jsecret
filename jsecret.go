package main

import (
	"bufio"
	"flag"
	"os"
	"sync"
)

var (
	HashList     = []string{}
	concurrency  int
)

func init() {
	flag.IntVar(&concurrency, "c", 5, "Concurrency level (number of goroutines)")
}

func main() {
	flag.Parse()
	var wg sync.WaitGroup
	scanner := bufio.NewScanner(os.Stdin)

	// Create a buffered channel to limit the number of concurrent goroutines
	semaphore := make(chan struct{}, concurrency)

	for scanner.Scan() {
		line := scanner.Text()
		if isUrl(line) {
			wg.Add(1)
			// Acquire a slot in the semaphore
			semaphore <- struct{}{}
			go func(url string) {
				defer wg.Done()
				matcher(url)
				// Release the slot after the goroutine finishes
				<-semaphore
			}(line)
		}
	}

	wg.Wait()
}

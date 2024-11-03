package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var HashList = []string{}
var maxConcurrentRequests = 5
var wg sync.WaitGroup

func main() {
	semaphore := make(chan struct{}, maxConcurrentRequests)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if isUrl(line) {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()
				matcher(url)
			}(line)
		}
	}

	wg.Wait()
}

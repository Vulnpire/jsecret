package main

import (
	"bufio"
	"flag"
	"os"
	"sync"
)

var HashList = []string{}

func main() {
	concurrency := flag.Int("c", 10, "number of concurrent threads")
	flag.Parse()

	var wg sync.WaitGroup
	scanner := bufio.NewScanner(os.Stdin)
	urls := make(chan string, *concurrency)

	
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range urls {
				if isUrl(url) {
					matcher(url)
				}
			}
		}()
	}

	for scanner.Scan() {
		urls <- scanner.Text()
	}
	close(urls)

	wg.Wait()
}

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
	regexFile    string
)

func main() {
	flag.IntVar(&concurrency, "c", 10, "Set the level of concurrency")
	flag.StringVar(&regexFile, "i", "", "Path to regex file")
	flag.Parse()

	// Load regex patterns from file
	loadRegexFromFile(regexFile)

	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if isUrl(line) {
			wg.Add(1)
			sem <- struct{}{} // Block if concurrency limit is reached
			go func(url string) {
				defer wg.Done()
				defer func() { <-sem }()
				matcher(url)
			}(line)
		}
	}

	wg.Wait()
	close(sem)
}

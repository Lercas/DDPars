package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/miekg/dns"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run parser.go <input_folder>")
		os.Exit(1)
	}

	inputFolder := os.Args[1]
	files, err := os.ReadDir(inputFolder)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	fileChan := make(chan string, 10)
	domainChan := make(chan string, 1000)
	done := make(chan struct{})

	go func() {
		outputFile, err := os.Create("domains.txt")
		if err != nil {
			fmt.Println("Error creating output file:", err)
			os.Exit(1)
		}
		defer outputFile.Close()
		domainSet := make(map[string]struct{})
		writer := bufio.NewWriter(outputFile)
		for domain := range domainChan {
			if _, exists := domainSet[domain]; !exists {
				domainSet[domain] = struct{}{}
				fmt.Fprintln(writer, domain)
			}
		}
		writer.Flush()
		close(done)
	}()

	
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range fileChan {
				parseFile(filePath, domainChan)
			}
		}()
	}

	for _, file := range files {
		if !file.IsDir() {
			fullPath := filepath.Join(inputFolder, file.Name())
			fileChan <- fullPath
		}
	}
	close(fileChan)

	wg.Wait()
	close(domainChan)
	<-done
}

func parseFile(filePath string, domainChan chan<- string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	parser := dns.NewZoneParser(file, "", "")
	for rr, ok := parser.Next(); ok; rr, ok = parser.Next() {
		if rr == nil {
			continue
		}
		domain := rr.Header().Name

		if strings.Contains(domain, "*") {
			continue
		}

		domain = strings.TrimSuffix(domain, ".")

		domain = strings.ReplaceAll(domain, "_", "")

		domainChan <- domain
	}
	if err := parser.Err(); err != nil && err != io.EOF {
		fmt.Println("Error parsing file:", err)
	}
}

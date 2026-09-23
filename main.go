package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type CheckResult struct {
	URL        string
	StatusCode int
	Duration   time.Duration
	Error      error
}

func checkURL(client *http.Client, url string) CheckResult {
	startedAt := time.Now()
	resp, err := client.Get(url)
	duration := time.Since(startedAt)

	if err != nil {
		return CheckResult{
			URL:      url,
			Duration: duration,
			Error:    err,
		}
	}

	defer resp.Body.Close()

	return CheckResult{
		URL:        url,
		StatusCode: resp.StatusCode,
		Duration:   duration,
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Ошибка: url ресурса не введен!")
		return
	}

	client := &http.Client{Timeout: 5 * time.Second}

	slots := make(chan int, 3)
	var wg sync.WaitGroup
	results := make(chan CheckResult, len(os.Args[1:]))

	for _, url := range os.Args[1:] {
		slots <- 1
		wg.Go(func() {
			defer func() {
				<-slots
			}()
			result := checkURL(client, url)
			results <- result
		})
	}

	wg.Wait()
	close(results)

	for result := range results {
		if result.Error != nil {
			fmt.Printf("Ошибка: %s | %s | %s \n", result.Error, result.URL, result.Duration)
		} else {
			fmt.Printf("%s | HTTP-статус: %d | %s \n", result.URL, result.StatusCode, result.Duration)
		}
	}

}

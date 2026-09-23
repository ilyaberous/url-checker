package main

import (
	"flag"
	"fmt"
	"net/http"
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
	concurrency := flag.Int("concurrency", 3, "Количество одновременных проверок")
	timeout := flag.Duration("timeout", 5*time.Second, "Таймаут запроса")
	flag.Parse()

	urls := flag.Args()

	if len(urls) == 0 {
		fmt.Println("Ошибка: url ресурсов не введен!")
		return
	}

	if *concurrency <= 0 {
		fmt.Println("Ошибка: concurrency должно быть больше нуля")
		return
	}

	if *timeout <= 0 {
		fmt.Println("Ошибка: timeout должен быть больше нуля")
		return
	}

	client := &http.Client{Timeout: *timeout}

	slots := make(chan int, *concurrency)
	var wg sync.WaitGroup
	results := make(chan CheckResult, len(urls))

	for _, url := range urls {
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

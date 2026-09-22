package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func checkURL(client *http.Client, url string) {
	startedAt := time.Now()

	resp, err := client.Get(url)

	fmt.Println("Время выполнения запроса к", url, ":", time.Since(startedAt))

	if err != nil {
		fmt.Println("Ошибка запроса к", url, ":", err)
		return
	}

	defer resp.Body.Close()

	fmt.Println("HTTP-статус запроса к", url, ":", resp.StatusCode)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Ошибка: url ресурса не введен!")
		return
	}

	client := &http.Client{Timeout: 5 * time.Second}

	for _, url := range os.Args[1:] {
		checkURL(client, url)
	}
}

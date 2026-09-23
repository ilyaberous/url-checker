package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckURLSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	))

	defer server.Close()

	result := checkURL(server.Client(), server.URL)

	//Проверка 1 - ответ должен быть без ошибки
	if result.Error != nil {
		t.Fatalf("Ожидался ответ без ошибки, но получен: %v", result.Error)
	}

	//Проверка 2 - статус код ответа должен быть OK
	if result.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус-код OK, но получен: %d", result.StatusCode)
	}

	// Проверка 3 — в результате сохранён переданный URL
	if result.URL != server.URL {
		t.Fatalf("Ожидался URL %s, получен %s", server.URL, result.URL)
	}
}
func TestCheckURLInvalidURL(t *testing.T) {
	client := &http.Client{}
	url := "://"

	result := checkURL(client, url)

	//Проверка 1 - должна быть ошибка = Error не должно быть nil
	if result.Error == nil {
		t.Fatalf("Ожидалась ошибка, но получен %v", result.Error)
	}

	//Проверка 2 - статус код ответа должен быть равен 0 = ответ не получен
	if result.StatusCode != 0 {
		t.Fatalf("В качестве кода ответа ожидался 0, а получен %d", result.StatusCode)
	}

	//Проверка 3 - в ответе сохранен переданный URL
	if result.URL != url {
		t.Fatalf("Ожидался URL %s, получен %s", url, result.URL)
	}
}

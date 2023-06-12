package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

func main() {
 endpoint := "http://localhost:8080/"
 fmt.Println("Введите длинный URL")
 var long string
 fmt.Scanln(&long)
 client := resty.New()
 response, err := client.R().
  SetHeader("Content-Type", "application/x-www-form-urlencoded").
  SetBody(strings.NewReader(fmt.Sprintf("url=%s", long))).
  Post(endpoint)
 if err != nil {
  fmt.Println("Ошибка при отправке запроса:", err)
  os.Exit(1)
 }
 fmt.Println("Статус-код", response.StatusCode())
 fmt.Println(response.String())
}
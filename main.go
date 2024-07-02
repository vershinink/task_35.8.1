package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"regexp"
	"time"
)

const (
	verbsAddr  string = "https://go-proverbs.github.io/"
	serverAddr string = "0.0.0.0:12345"
	proto      string = "tcp4"
)

func main() {
	// Получаем список Go поговорок.
	proverbs, err := proverbs(verbsAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Запускаем сетевую службу по протоколу TCP на порту 12345.
	listener, err := net.Listen(proto, serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// Запускаем прием входящих соединений в цикле.
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go verbsHandle(conn, proverbs)
	}
}

// proverbs - получает список Go поговорок с переданного url
// и возвращает его и ошибку в виде двумерного слайса байт.
func proverbs(path string) ([][]byte, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to proverbs website: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body: %w", err)
	}

	re := regexp.MustCompile(`<h3><a href=".*">(.*)</a></h3>`)
	var verbs [][]byte

	sub := re.FindAllSubmatch(body, -1)
	line := []byte("\r\n")
	for _, s := range sub {
		s[1] = append(s[1], line...)
		verbs = append(verbs, s[1])
	}
	return verbs, nil
}

// verbsHandle - обработчик входящего подключения. Раз в 3 секунды передает
// клиенту случайную Go поговорку из списка verbs.
func verbsHandle(conn net.Conn, verbs [][]byte) {
	for {
		i := rand.Intn(len(verbs))
		conn.Write(verbs[i])
		time.Sleep(3 * time.Second)
	}
}

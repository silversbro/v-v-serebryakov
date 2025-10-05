package main

import (
	"context"
	"log"
	"net"
	"time"
)

const greetingMessage = "Добро пожаловать на сервер\n Сегодня хорошая погода \n и настроение отличное"

func main() {
	localAddr := "0.0.0.0:52521"
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Fatalf("Соединение открыть не удалось, видимо порт %s занят", localAddr)
	}
	log.Printf("Успешно начато прослушивание порта %s для входящих соединений", localAddr)

	ctx, _ := context.WithTimeout(context.Background(), time.Minute*15)

	go func() {
		select {
		case <-ctx.Done():
			log.Print("Пришло время выключать сервер")
			listener.Close()
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Этот листенер сломался, несите следующий")
		}
		log.Printf("Установлено входящее соединение с адреса %s", conn.RemoteAddr())
		conn.Write([]byte(greetingMessage))
		time.Sleep(3 * time.Second)
		conn.Write([]byte("\n"))
	}
	//time.Sleep(30 * time.Second)

	defer listener.Close()

}

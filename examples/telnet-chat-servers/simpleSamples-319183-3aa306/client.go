package main

import (
	"bufio"
	"log"
	"net"
	"time"
)

func main() {
	remoteAddr := "127.0.0.1:52521"
	conn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Fatalf("Не удалось открыть соединение, по причине: %v", err)
	}
	log.Printf("Удалось открыть соединение с адреса %s на адрес %s", conn.LocalAddr(), conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Fatalf("Что-то пошло не так при чтении приветственного сообщения от сервера. %v", err)
		} else {
			log.Printf("Сервер нам прислал: %s", scanner.Text())
		}
	}
	log.Print("Соединение разорвано")

	time.Sleep(time.Minute * 15)

}

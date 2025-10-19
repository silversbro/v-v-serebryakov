package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		log.Printf("Usage: %s [--timeout=10s] host port\n", os.Args[0])
		os.Exit(1)
	}

	host, port := args[0], args[1]
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Printf("Error connecting to %s: %v\n", address, err)
		os.Exit(1)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Error closing client: %v", err)
		}
	}()

	log.Printf("...Connected to %s\n", address)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	done := make(chan struct{}, 2) // буфер для 2 горутин

	go func() {
		if err := client.Send(); err != nil {
			handleError(err, "send")
		}
		done <- struct{}{}
	}()

	go func() {
		if err := client.Receive(); err != nil {
			handleError(err, "receive")
		}
		done <- struct{}{}
	}()

	select {
	case <-sigCh:
		log.Println("...SIGINT received, closing connection")
	case <-done:
		log.Println("...Connection closed")
	}
}

func handleError(err error, operation string) {
	if errors.Is(err, io.EOF) {
		if operation == "send" {
			log.Println("...EOF")
		} else {
			log.Println("...Connection was closed by peer")
		}
	} else {
		var opErr *net.OpError
		if errors.As(err, &opErr) && opErr.Op == "read" {
			log.Println("...Connection was closed by peer")
		} else {
			log.Printf("Error in %s: %v", operation, err)
		}
	}
}

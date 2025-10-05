package main

import (
	"flag"
	"fmt"
	"io"
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
		fmt.Println(os.Stderr, "Usage: %s [--timeout=10s] host port\n", os.Args[0])
		os.Exit(1)
	}

	host, port := args[0], args[1]
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Println("Error connecting to %s: %v\n", address, err)
		os.Exit(1)
	}
	defer client.Close()

	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	done := make(chan struct{})
	go runClient(client, done)

	select {
	case <-sigCh:
		fmt.Println(os.Stderr, "...SIGINT received, closing connection\n")
		return
	case <-done:
		fmt.Println(os.Stderr, " closing connection\n")
		return
	}
}

func runClient(client TelnetClient, done chan struct{}) {
	defer close(done)

	sendErr := make(chan error, 1)
	receiveErr := make(chan error, 1)

	go func() {
		if err := client.Send(); err != nil {
			sendErr <- err
		}
	}()

	go func() {
		if err := client.Receive(); err != nil {
			receiveErr <- err
		}
	}()

	select {
	case err := <-sendErr:
		handleError(err, "send")
	case err := <-receiveErr:
		handleError(err, "receive")
	}
}

func handleError(err error, operation string) {
	if err == io.EOF {
		if operation == "send" {
			fmt.Println(os.Stderr, "...EOF\n")
		} else {
			fmt.Println(os.Stderr, "...Connection was closed by peer\n")
		}
	} else if _, ok := err.(*net.OpError); ok {
		fmt.Println(os.Stderr, "...Connection was closed by peer\n")
	} else {
		fmt.Println(os.Stderr, "Error in %s: %v\n", operation, err)
	}
}

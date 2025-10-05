package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func readRoutine(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(conn)
	defer log.Printf("Finished readRoutine")

OUTER:
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !scanner.Scan() {
				log.Printf("CANNOT SCAN")
				break OUTER
			}
			text := scanner.Text()
			log.Printf("From server: %s", text)
		}
	}
}

func writeRoutine(ctx context.Context, conn net.Conn, wg *sync.WaitGroup, stdin chan string) {
	defer wg.Done()
	defer log.Printf("Finished writeRoutine")

	for {
		select {
		case <-ctx.Done():
			return
		case str := <-stdin:
			log.Printf("To server %v\n", str)

			conn.Write([]byte(fmt.Sprintf("%s\n", str)))
		}

	}
}

func stdinScan() chan string {
	out := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			out <- scanner.Text()
		}
		if scanner.Err() != nil {
			close(out)
		}
	}()
	return out
}

func main() {
	dialer := &net.Dialer{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := dialer.DialContext(ctx, "tcp", "127.0.0.1:3302")
	defer conn.Close()
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}

	stdin := stdinScan()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		readRoutine(ctx, conn, wg)
		cancel()
	}()

	wg.Add(1)
	go func() {
		writeRoutine(ctx, conn, wg, stdin)
	}()

	wg.Wait()
}

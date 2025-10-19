package main

import (
	"bufio"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
	scanner *bufio.Scanner
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}
	t.conn = conn
	t.scanner = bufio.NewScanner(conn)
	return nil
}

func (t *telnetClient) Close() error {
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}

func (t *telnetClient) Send() error {
	scanner := bufio.NewScanner(t.in)
	for scanner.Scan() {
		data := scanner.Bytes()
		data = append(data, '\n')
		if _, err := t.conn.Write(data); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func (t *telnetClient) Receive() error {
	if t.scanner == nil {
		return nil
	}

	for t.scanner.Scan() {
		data := t.scanner.Bytes()
		data = append(data, '\n')
		if _, err := t.out.Write(data); err != nil {
			return err
		}
	}
	return t.scanner.Err()
}

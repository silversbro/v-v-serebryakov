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
	return nil
}

func (t *telnetClient) Close() error {
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}

func (t *telnetClient) Send() error {
	reader := bufio.NewReader(t.in)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		if _, err := t.conn.Write([]byte(line)); err != nil {
			return err
		}
	}
}

func (t *telnetClient) Receive() error {
	reader := bufio.NewReader(t.conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		if _, err := t.out.Write([]byte(line)); err != nil {
			return err
		}
	}
}
